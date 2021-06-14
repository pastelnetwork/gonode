import argparse
import os
import shutil
import time
import requests
from bs4 import BeautifulSoup
from selenium import webdriver
from time import sleep
import cairosvg
import json
import re

# GLOBAL variables / switches
# Selenium's driver
# Currently set to chrome, supposing
# Chrome and chromdriver is installed
option = webdriver.ChromeOptions()
chrome_prefs = {}
option.experimental_options["prefs"] = chrome_prefs
chrome_prefs["profile.default_content_settings"] = {"images": 2}
chrome_prefs["profile.managed_default_content_settings"] = {"images": 2}

CATEGORY = 'art'

DRIVER = webdriver.Chrome()
DRIVER.maximize_window()
# URL to the webpage
BASE_URL = 'https://opensea.io'
# Pause time in seconds
SCROLL_PAUSE_TIME = 0.4
# In case it is set to 1 it is infinite !!!
SCROLLING_ITERATION = 3 # Let's be it 10000 or end of page
# Storing URLs dictionary for fast lookup of keys, so URLs = Keys
NORMAL_SIZE_IMG_URLS = {}
# Calling script function purely (here or
# from another file) or via command line
USE_CMD_ARGS = 0
# Input - output folders
SAVE_FOLDER = ''
# NFTs Collections Names
ASSET_COLLECTION_NAME = []

LAST_FOUND = False

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
        SAVE_FOLDER = "scraper_categories/"


def retrieve_title(original_title):
    # Replacing special characters / strings, and numbers
    #title = original_title.replace(' ', '_')
    #title = title.replace('OpenSea', '')

    title = original_title.replace('https://', '')
    title = original_title.replace('/', '_')

    #regex = re.compile('[^a-zA-Z]')
    # First parameter is the replacement, second parameter is your input string
    #title = regex.sub('', title)

    return title


def parse_data(end_calculation, collection):

    global DRIVER
    html = DRIVER.page_source

    # BeautifulSoup obj for scraping static data from html
    soup = BeautifulSoup(html, 'html.parser')
    find_all_a = soup.find_all("a", href=True)

    # Iterate through the href-s of thumbnails
    # and determine the URLs which contains an /asset (=image)
    for el in find_all_a:
        if el['href'].startswith('/assets/'):
            url = BASE_URL + el['href']

            # Saving it into a dictionary due to speed
            # of lookup and insert
            if url not in NORMAL_SIZE_IMG_URLS:
                end_calculation = 0
                print("URL appended: {}".format(url))
                NORMAL_SIZE_IMG_URLS[url] = collection


def download_images():
    global NORMAL_SIZE_IMG_URLS

    # Iterate through keys of the dictionary
    # which are the URLs
    for url, collection in NORMAL_SIZE_IMG_URLS.items():

        try:
            normal_img_page = requests.get(url)
            soup2 = BeautifulSoup(normal_img_page.content, 'html.parser')

            title = soup2.find("meta", property="og:title")['content']
            img_url = soup2.find("meta", property="og:image")['content']

            # save image
            print("Savingimage")
            to_be_saved_name = collection + '-' + retrieve_title(url)
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
            elif file_type.endswith('svg+xml'):
                to_be_saved_name = to_be_saved_name + '.svg'
            else:
                print("Can't save image of file_type " + file_type)
                del downloaded_img
                continue

            to_be_saved_name = os.path.join(SAVE_FOLDER+'{}/'.format(CATEGORY), to_be_saved_name)
            print("To be saved to: {}".format(to_be_saved_name))
            with open(to_be_saved_name, 'wb') as out_file:
                shutil.copyfileobj(downloaded_img.raw, out_file)

            if file_type.endswith('svg+xml'):
                print("Converting from svg to png: {}".format(to_be_saved_name))
                cairosvg.svg2png(url=to_be_saved_name, write_to=os.path.join(SAVE_FOLDER+'{}/'.format(CATEGORY), collection + '-' + retrieve_title(url)+".png"))
                os.remove(to_be_saved_name)

            del downloaded_img
        except Exception as bs:
            print ("Base exception is in saving: {}".format(bs))
    pass

def fetch_urls_and_images():
    global NORMAL_SIZE_IMG_URLS
    global DRIVER
    global ASSET_COLLECTION_NAME

    print(str(ASSET_COLLECTION_NAME))

    for collection in ASSET_COLLECTION_NAME:
        url = 'https://opensea.io/assets/{}'.format(collection)
        DRIVER.get(url)
        # Allow 2 seconds for the page to be opened
        sleep(3)
        html = DRIVER.page_source
        sleep(SCROLL_PAUSE_TIME)

        # Get the screen height
        screen_height = DRIVER.execute_script("return window.screen.height")
        i = 1
        j = 1

        # To check if we reached the end
        end_calculation = 0

        while True:
            # Scroll one screen height each time
            DRIVER.execute_script(f"window.scrollTo(0, {screen_height}*{i});".format(screen_height=screen_height, i=i))
            i += 1
            j += 1

            time.sleep(SCROLL_PAUSE_TIME)
            # Update scroll height each time after scrolled
            # as this will be changed each and every iteration
            scroll_height = DRIVER.execute_script("return document.body.scrollHeight;")
            print ("scroll_height is: {}".format(scroll_height))
            end_calculation += 1

            parse_data(end_calculation, collection)

            if (end_calculation >= 320):
                # It means we have reached the end of the category
                print("END PAGE")
                download_images()
                NORMAL_SIZE_IMG_URLS.clear()
                i = SCROLLING_ITERATION+1

            print("len(NORMAL_SIZE_IMG_URLS)="+str(len(NORMAL_SIZE_IMG_URLS)))

            if(len(NORMAL_SIZE_IMG_URLS) > 2000):
                end_calculation = 0
                print ("Scrolling iteration is : {}".format(i))
                print("DOWNLOAD")
                download_images()
                NORMAL_SIZE_IMG_URLS.clear()
                print("Screen height is: {}".format(screen_height))


            if (j > 150 ):
                j = 0
                end_calculation = 0
                # We need to re-init driver
                print("Driver REINIT to save time")
                DRIVER.quit()

                DRIVER = webdriver.Chrome()
                DRIVER.maximize_window()

                DRIVER.get(url)
                # Allow 2 seconds for the page to be opened
                sleep(3)
                html = DRIVER.page_source
                sleep(5)

                DRIVER.execute_script(f"window.scrollTo(0, {screen_height}*{i});".format(screen_height=screen_height, i=i))

                time.sleep(10)
                print("waiting due to reinitialization")
                # Update scroll height each time after scrolled
                # as this will be changed each and every iteration
                #scroll_height = DRIVER.execute_script("return document.body.scrollHeight;")
                #i += 1



            # If SCROLLING_ITERATION is set to 1 it will also apply infinite scroll
            if i > SCROLLING_ITERATION:# or screen_height * i > scroll_height --> nem jo mert nem adja vissza a dolgokat
                #print("Params:\ni:{}\nscreen_height:{}\nscroll_height:{}".format(i,screen))
                i = 0
                j = 0
                end_calculation = 0
                break

    download_images()

def fetch_collections_names():
    global ASSET_COLLECTION_NAME
    global DRIVER

    url = 'https://opensea.io/assets?search[categories][0]={}'.format(CATEGORY)
    DRIVER.get(url)
    # Allow 2 seconds for the page to be opened
    sleep(1)
    html = DRIVER.page_source
    sleep(SCROLL_PAUSE_TIME)

    while True:
        scroll_height = DRIVER.execute_script('return document.getElementsByClassName("Scrollbox--content")[0].scrollHeight;')
        initial_collections_len = len(ASSET_COLLECTION_NAME)

        # BeautifulSoup obj for scraping static data from html
        soup = BeautifulSoup(html, 'html.parser')
        next_data_script = soup.find('script', attrs={ 'id': '__NEXT_DATA__'})
        content_string = str(next_data_script.contents)
        collections_slugs = re.findall('\"slug\"\:\"([a-zA-Z0-9-]+)\"', content_string)

        for slug in collections_slugs:
            if slug not in NORMAL_SIZE_IMG_URLS:
                ASSET_COLLECTION_NAME.append(slug)
        
        print('ASSET_COLLECTION_NAME='+str(len(ASSET_COLLECTION_NAME)))
        if initial_collections_len == len(ASSET_COLLECTION_NAME):
            break

        DRIVER.execute_script('document.getElementsByClassName("Scrollbox--content")[0].scrollTop = {};'.format(scroll_height))
        time.sleep(2)

        if scroll_height >= 105000:
            break

if __name__ == '__main__':
    parse_folders()
    try:
        fetch_collections_names()
        fetch_urls_and_images()
    except Exception as base_ex:
        print(base_ex)
