from selenium import webdriver
from selenium.webdriver.common.by import By
from selenium.webdriver.common.keys import Keys
from selenium.webdriver.chrome.service import Service
from selenium.webdriver.support.ui import WebDriverWait
from selenium.webdriver.support import expected_conditions as EC
from webdriver_manager.chrome import ChromeDriverManager
from selenium.common.exceptions import NoSuchElementException, StaleElementReferenceException

import requests
import time

# Настройка Selenium WebDriver
driver = webdriver.Chrome()
wait = WebDriverWait(driver, 10)

# Базовый URL
BASE_URL = "http://localhost:7540"

# API endpoints
LOGIN_API = f"{BASE_URL}/api/signin"
CREATE_TASK_API = f"{BASE_URL}/api/task"

# Test data
TEST_USER_PASSWORD = "123456789"
TEST_TASK = {
  "comment": "comment",
  "date": "31.12.2024",
  "id": "123456788",
  "repeat": "d 12",
  "title": "qwerty"
}

try:

    # Открытие базового URL
    driver.get(BASE_URL)

    # Ожидание загрузки страницы
    wait.until(EC.presence_of_element_located((By.TAG_NAME, "body")))

    # STEP 1: Авторизация через API
    print("Авторизация через API...")
    login_response = requests.post(LOGIN_API, json={"password": "123456789"})



    # Получение токена из ответа
    token = login_response.json().get("token")
    if not token:
        raise Exception("Токен авторизации не получен.")

    headers = {"Authorization": f"Bearer {token}"}
    print(headers)

    # STEP 2: Создание задачи через API
    print("Создание задачи через API...")
    add_task_button = driver.find_element(By.CSS_SELECTOR, ".btn.primary")


    assert add_task_button.is_displayed(), "Кнопка 'Добавить задачу' не отображается"
    assert add_task_button.is_enabled(), "Кнопка 'Добавить задачу' не активна"

    add_task_button.click()

    time.sleep(2)

    modal_window = driver.find_element(By.CSS_SELECTOR, ".modal .dialog")
    assert modal_window.is_displayed(), "Модальное окно 'Добавить новую задачу' не отображается"

    # STEP 4: Заполнение формы
    print("Заполнение формы...")
    # Заполнение поля "Дата"
    date_input = modal_window.find_element(By.CSS_SELECTOR,
                                           'input.no-right-radius')  # Укажите правильный placeholder
    date_input.clear()
    date_input.send_keys(TEST_TASK["date"])

    # Заполнение поля "Задача"
    task_input = modal_window.find_element(By.CSS_SELECTOR,
                                           ".form-input:nth-child(2) .input")  # Укажите правильный placeholder
    task_input.clear()
    task_input.send_keys(TEST_TASK["title"])

    # Заполнение поля "Комментарий"
    comment_input = modal_window.find_element(By.CSS_SELECTOR,
                                              ".form-input:nth-child(3) .input")  # Укажите правильный placeholder
    comment_input.clear()
    comment_input.send_keys(TEST_TASK["comment"])

    # STEP 5: Подтверждение добавления задачи
    print("Подтверждение добавления задачи...")
    submit_button = modal_window.find_element(By.CSS_SELECTOR,
                                              'button.btn.primary')  # Укажите правильный селектор кнопки
    assert submit_button.is_displayed() and submit_button.is_enabled(), "Кнопка 'Добавить' не отображается или не активна"
    submit_button.click()

    # Ожидание закрытия модального окна
    wait.until(EC.invisibility_of_element(modal_window))
    print("Задача успешно добавлена.")

    # ожидаем исчезновения модального окна
    WebDriverWait(driver, 10).until(
        EC.invisibility_of_element_located((By.CSS_SELECTOR, ".dialog"))
    )

    # проверяем, что задача добавлена
    task_list = driver.find_element(By.CSS_SELECTOR, ".notelist")
    tasks = task_list.find_elements(By.CSS_SELECTOR, ".notecard")
    assert len(tasks) == 1

    # STEP 6: Удаление задачи
    print("Удаление задачи...")
    time.sleep(1)  # Ждем загрузки списка задач (если требуется)
    done_task_btn = driver.find_element(By.CSS_SELECTOR, ".btnicon")
    done_task_btn.click()

    # Проверка, что список задач пуст
    try:
        task_list = driver.find_element(By.CSS_SELECTOR, ".notelist")
        tasks = task_list.find_elements(By.CSS_SELECTOR, ".notecard")
        assert len(tasks) == 0, f"Ожидалось 0 задач, найдено: {len(tasks)}"
    except (NoSuchElementException, StaleElementReferenceException):
        print("Список задач пуст — тест успешно завершен.")


finally:
    # Завершение работы
    driver.quit()

