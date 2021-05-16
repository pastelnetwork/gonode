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
wget https://storage.googleapis.com/tensorflow/libtensorflow/libtensorflow-cpu-linux-x86_64-2.4.0.tar.gz
sudo tar -C /usr/local -xzf ./libtensorflow-cpu-linux-x86_64-2.4.0.tar.gz
sudo /sbin/ldconfig -v
```

Download ML models and images:

```
pip install gdown
~/.local/bin/gdown https://drive.google.com/uc?id=1U6tpIpZBxqxIyFej2EeQ-SbLcO_lVNfu
unzip ./SavedMLModels.zip -d ./

wget https://www.dropbox.com/s/6ohzgvz418rhl4l/Animecoin_All_Finished_Works.zip
unzip Animecoin_All_Finished_Works.zip -d ./allRegisteredWorks

wget https://www.dropbox.com/s/4uajzyh09bc0rp3/dupe_detector_test_images.zip
unzip dupe_detector_test_images.zip -d ./dupes

wget https://www.dropbox.com/s/yjqsxsz97msai4e/non_duplicate_test_images.zip
unzip non_duplicate_test_images.zip -d ./originals
```

Download the latest test corpus of images:
```
~/.local/bin/gdown https://drive.google.com/uc?id=1BslINgdqs8ik7PiDjKKL1wlrRfQanVjQ
unzip test_corpus_opens_1.zip -d ./test_corpus
```

## Goptuna Optimizer

[Optimizer application](./cmd/optimizer) is configured with cmaes sampler to find the maximum AUPRC.

# Pre-requisites

`goptuna` cli tool executable should be downloaded from [the project's GitHub Releases](https://github.com/c-bata/goptuna/releases) page.

# Workflow

Goptuna studies results are saved into MySQL database.

Default setup would work with Docker MySQL images:

```
docker pull mysql:8.0

export GOPTUNA_CONTAINER=$(docker run   -d   --rm   -p 3306:3306   -e MYSQL_USER=goptuna   -e MYSQL_DATABASE=goptuna   -e MYSQL_PASSWORD=password   -e MYSQL_ALLOW_EMPTY_PASSWORD=yes   --name goptuna-mysql   mysql:8.0)
```

Create initial goptuna database structure and empty study:

```
./goptuna create-study --storage mysql://goptuna:password@localhost:3306/goptuna --study dupe-detection-aurpc
```

Run goptuna dashboard to observe studies results:
```
./goptuna dashboard --storage mysql://goptuna:password@127.0.0.1:3306/goptuna
```

Backup goptuna database data before it is vanished with termination of docker container:

```
docker exec $GOPTUNA_CONTAINER sh -c 'exec mysqldump --no-tablespaces --databases goptuna -ugoptuna -ppassword' > backup.sql
```

Restore goptuna database data from backup:
```
docker exec -i $GOPTUNA_CONTAINER sh -c 'exec mysql -ugoptuna -ppassword' < ./backup.sql
```

Run optimizer with `imageCount` parameter to limit number of analyzed images per trial;  
`runCount` defines the number of trials per run;  
`studyName` defines the name of Goptuna study;  
`rootDir` defines the directory from where to load the corpus of images and where during the first run to generate sqlite database with fingerprints.
```
go run ./cmd/optimizer/ -rootDir "./test_corpus" -imageCount 30 -runCount 100 -studyName "dupe-detection-aurpc"
```