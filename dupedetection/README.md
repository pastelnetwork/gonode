[![Maintained by PastelNetwork](https://img.shields.io/badge/maintained%20by-pastel.network-%235849a6.svg)](https://pastel.network)

# Dupe detection API

This repo expose an API for internal use.
The API receives input-image and return fingerprints, rareness scores and some other info.
The API uses two parameters from SuperNode config file:
```
input_dir: dupe detection service input directory path
output_dir: dupe detection service output directory path
```

The API will:
1. stores input image into the dd-service input directory  using some random number as file name.
2. waits for the output JSON file with the same name to appear in the dd-service output directory.
3. parses JSON output and return result to the caller.

## To test the API on local

### Set up environment to run dd-service (Ubuntu 20.04)

First install required libraries:    
```
pip install xgboost hyppo zstandard tensorflow pandas scipy scikit-learn matplotlib watchdog chromedriver_autoinstaller selenium Pillow opennsfw-standalone tensorflow_hub imagehash
```
Then install Google Chrome (needed for "rare on the internet" code)
```bash
wget https://dl.google.com/linux/direct/google-chrome-stable_current_amd64.deb
sudo apt install ./google-chrome-stable_current_amd64.deb
```
Then create the following folder structure:

```
mkdir -p /home/$USER/pastel_dupe_detection_service/input_files/
mkdir -p /home/$USER/pastel_dupe_detection_service/support_files/
mkdir -p /home/$USER/pastel_dupe_detection_service/output_files/
mkdir -p /home/$USER/pastel_dupe_detection_service/processed_files/
mkdir -p /home/$USER/pastel_dupe_detection_service/rare_on_internet/
mkdir -p /home/$USER/pastel_dupe_detection_service/support_files/mobilenet_v2_140_224/
```

Then, in the "support_files" directory, put the following files:

* dupe_detection_image_fingerprint_database.sqlite (download and extract from: https://download.pastel.network/machine-learning/registered_image_fingerprints_db.sqlite )

* DupeDetector_gray.pth.tar (download from: https://download.pastel.network/machine-learning/DupeDetector_gray.pth.tar)

* pca_bw.vt (download from: https://download.pastel.network/machine-learning/pca_bw.vt)

* train_0_bw.hdf5 (download from: https://download.pastel.network/machine-learning/train_0_bw.hdf5)

* config.ini
```ini
[DUPEDETECTIONCONFIG]
input_files_path = /home/$USER/pastel_dupe_detection_service/input_files/
support_files_path = /home/$USER/pastel_dupe_detection_service/support_files/
output_files_path = /home/$USER/pastel_dupe_detection_service/output_files/
processed_files_path = /home/$USER/pastel_dupe_detection_service/processed_files/
internet_rareness_downloaded_images_path = /home/$USER/pastel_dupe_detection_service/rare_on_internet/
nsfw_model_path = /home/$USER/pastel_dupe_detection_service/support_files/mobilenet_v2_140_224/

```

Then, in the "mobilenet_v2_140_224" directory, put the following file:
saved model (download and extract from: https://download.pastel.network/machine-learning/nsfw_mobilenet_v2_140_224.zip)

Then create the following environmental variable to store the config.ini path:
```
export DUPEDETECTIONCONFIGPATH=/home/$USER/pastel_dupe_detection_service/support_files/config.ini
```

Clone dd-service source: https://github.com/pastelnetwork/dd-service

### Run test

#### Run dd-service:
```bash
cd <path to dd-service>/dd-service
python3 dupe_detection_server.py
```

#### Run test API:
```bash
cd cmd
go run test.go -workDir /home/$USER/pastel_dupe_detection_service/
```
it takes about 2 minutes for the test end.
