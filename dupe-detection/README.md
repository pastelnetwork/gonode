# dupe-detection

Used ML models can be downloaded from [there](https://drive.google.com/file/d/1U6tpIpZBxqxIyFej2EeQ-SbLcO_lVNfu/view?usp=sharing).

## Hardware requirements:

- at least 8GB of RAM

## Ubuntu

Install `swig`

```
sudo apt-get install -y swig
```

Install [relevant tensoflow C library](https://www.tensorflow.org/install/lang_c):

```
wget https://storage.googleapis.com/tensorflow/libtensorflow/libtensorflow-cpu-linux-x86_64-2.4.0.tar.g
sudo tar -C /usr/local -xzf ./libtensorflow-cpu-linux-x86_64-2.4.0.tar.gz
sudo /sbin/ldconfig -v
```

Download ML models and images:

```
pip install gdown
~/.local/bin/gdown https://drive.google.com/uc?id=1U6tpIpZBxqxIyFej2EeQ-SbLcO_lVNfu
unzip ./SavedMLModels.zip -d ./

wget https://www.dropbox.com/s/6ohzgvz418rhl4l/Animecoin_All_Finished_Works.zip
unzip Animecoin_All_Finished_Works.zip

wget https://www.dropbox.com/s/4uajzyh09bc0rp3/dupe_detector_test_images.zip
unzip dupe_detector_test_images.zip

wget https://www.dropbox.com/s/yjqsxsz97msai4e/non_duplicate_test_images.zip
unzip non_duplicate_test_images.zip
```