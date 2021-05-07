import argparse
import os
import shutil
import time
import re
import requests
from bs4 import BeautifulSoup
from selenium import webdriver
from time import sleep

# GLOBAL variables / switches
# Selenium's driver
# Currently set to chrome, supposing
# Chrome and chromdriver is installed
DRIVER = webdriver.Chrome()
DRIVER.maximize_window()
# URL to the webpage
BASE_URL = 'https://opensea.io'
# Pause time in seconds
SCROLL_PAUSE_TIME = 1
# In case it is set to 1 it is infinite !!!
SCROLLING_ITERATION = 10
# Storing URLs dictionary for fast lookup of keys, so URLs = Keys
NORMAL_SIZE_IMG_URLS = {}
# Calling script function purely (here or
# from another file) or via command line
USE_CMD_ARGS = 0
# Input - output folders
SAVE_FOLDER = ''


def parse_folders():
    global SAVE_FOLDER

    if USE_CMD_ARGS:
        # Construct the argument parse and parse the arguments
        ap = argparse.ArgumentParser()
        ap.add_argument("-i", "--inputFolder", type=str, required=True,
                        help="path to optional input folder")
        args = vars(ap.parse_args())

        SAVE_FOLDER = args["inputFolder"]
    else:
        SAVE_FOLDER = "input_folder/"


def retrieve_title(original_title):
    # Replacing special characters / strings, and numbers
    title = original_title.replace(' ', '_')
    title = title.replace('OpenSea', '')

    regex = re.compile('[^a-zA-Z]')
    # First parameter is the replacement, second parameter is your input string
    title = regex.sub('', title)

    return title


def parse_data():

    global DRIVER
    html = DRIVER.page_source

    # BeautifulSoup obj for scraping static data from html
    soup = BeautifulSoup(html, 'html.parser')
    find_all_a = soup.find_all("a", href=True)

    # Iterate through the href-s of thumbnails
    # and determine the URLs which contains an /asset (=image)
    for el in find_all_a:
        if el['href'].startswith('/assets/'):
            dict_key = BASE_URL + el['href']

            # Saving it into a dictionary due to speed
            # of lookup and insert
            if dict_key not in NORMAL_SIZE_IMG_URLS:
                print("URL appended: {}".format(dict_key))
                NORMAL_SIZE_IMG_URLS[dict_key] = 1


def download_images():
    global NORMAL_SIZE_IMG_URLS

    # Iterate through keys of the dictionary
    # which are the URLs
    for url in NORMAL_SIZE_IMG_URLS.keys():
        normal_img_page = requests.get(url)
        soup2 = BeautifulSoup(normal_img_page.content, 'html.parser')

        title = soup2.find("meta", property="og:title")['content']
        img_url = soup2.find("meta", property="og:image")['content']

        to_be_saved_name = retrieve_title(title)
        print(img_url)
        print("Saving image: {}".format(to_be_saved_name))

        # Download image and based on header's
        # "Content-type" determine the extension
        downloaded_img = requests.get(img_url, stream=True)
        file_type = downloaded_img.headers['Content-Type']
        if file_type.endswith('png'):
            to_be_saved_name = to_be_saved_name + '.png'
        elif file_type.endswith('jpg') or file_type.endswith('jpeg'):
            to_be_saved_name = to_be_saved_name + '.jpg'
        else:
            del downloaded_img
            continue

        to_be_saved_name = os.path.join(SAVE_FOLDER, to_be_saved_name)
        with open(to_be_saved_name, 'wb') as out_file:
            shutil.copyfileobj(downloaded_img.raw, out_file)

        del downloaded_img
    pass


def fetch_urls_and_images():
    global DRIVER
    url = 'https://opensea.io/assets'
    DRIVER.get(url)
    # Allow 2 seconds for the page to be opened
    sleep(2)
    html = DRIVER.page_source
    sleep(SCROLL_PAUSE_TIME)

    # Get the screen height
    screen_height = DRIVER.execute_script("return window.screen.height")
    i = 1

    while True:
        # Scroll one screen height each time
        DRIVER.execute_script(f"window.scrollTo(0, {screen_height}*{i});".format(screen_height=screen_height, i=i))
        i += 1
        time.sleep(SCROLL_PAUSE_TIME)

        # Update scroll height each time after scrolled
        # as this will be changed each and every iteration
        scroll_height = DRIVER.execute_script("return document.body.scrollHeight;")

        # Parse the currently available static html data
        parse_data()

        # If SCROLLING_ITERATION is set to 1 it will also apply infinite scroll
        if screen_height * i > scroll_height or i == SCROLLING_ITERATION:
            break;

    download_images()


if __name__ == '__main__':
    parse_folders()
    try:
        fetch_urls_and_images()
    except Exception as base_ex:
        print(base_ex)
