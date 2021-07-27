#!/usr/bin/python3
import os, tarfile, subprocess


#USER DEFINED
PASTELD_DIR = '/home/nd/pastelnetwork/pastel/src'

#No need to modify
ARCHIVE = './nodes.tar.gz'
EXTRACT_FOLDER = 'unzipped'
EXTRACT_PATH_FULL= os.getcwd() + "/" + EXTRACT_FOLDER


CMD_LIST=[
"{}/pasteld -datadir={}/node0  -daemon -keypool=1 -discover=0 -rest -debug=masternode,mnpayments,governance -masternode -txindex=1 -reindex -masternodeprivkey=91sY9h4AQ62bAhNk1aJ7uJeSnQzSFtz7QmW5imrKmiACm7QJLXe".format(PASTELD_DIR, EXTRACT_PATH_FULL),
"{}/pasteld -datadir={}/node1  -daemon -keypool=1 -discover=0 -rest -debug=masternode,mnpayments,governance -masternode -txindex=1 -reindex -masternodeprivkey=923JtwGJqK6mwmzVkLiG6mbLkhk1ofKE1addiM8CYpCHFdHDNGo".format(PASTELD_DIR, EXTRACT_PATH_FULL),
"{}/pasteld -datadir={}/node2  -daemon -keypool=1 -discover=0 -rest -debug=masternode,mnpayments,governance -masternode -txindex=1 -reindex -masternodeprivkey=91wLgtFJxdSRLJGTtbzns5YQYFtyYLwHhqgj19qnrLCa1j5Hp5Z".format(PASTELD_DIR, EXTRACT_PATH_FULL),
"{}/pasteld -datadir={}/node3  -daemon -keypool=1 -discover=0 -rest -debug=masternode,mnpayments,governance -masternode -txindex=1 -reindex -masternodeprivkey=92XctTrjQbRwEAAMNEwKqbiSAJsBNuiR2B8vhkzDX4ZWQXrckZv".format(PASTELD_DIR, EXTRACT_PATH_FULL),
"{}/pasteld -datadir={}/node4  -daemon -keypool=1 -discover=0 -rest -debug=masternode,mnpayments,governance -masternode -txindex=1 -reindex -masternodeprivkey=923JCnYet1pNehN6Dy4Ddta1cXnmpSiZSLbtB9sMRM1r85TWym6".format(PASTELD_DIR, EXTRACT_PATH_FULL),
"{}/pasteld -datadir={}/node5  -daemon -keypool=1 -discover=0 -rest -debug=masternode,mnpayments,governance -masternode -txindex=1 -reindex -masternodeprivkey=93BdbmxmGp6EtxFEX17FNqs2rQfLD5FMPWoN1W58KEQR24p8A6j".format(PASTELD_DIR, EXTRACT_PATH_FULL),
"{}/pasteld -datadir={}/node6  -daemon -keypool=1 -discover=0 -rest -debug=masternode,mnpayments,governance -masternode -txindex=1 -reindex -masternodeprivkey=92av9uhRBgwv5ugeNDrioyDJ6TADrM6SP7xoEqGMnRPn25nzviq".format(PASTELD_DIR, EXTRACT_PATH_FULL),
"{}/pasteld -datadir={}/node7  -daemon -keypool=1 -discover=0 -rest -debug=masternode,mnpayments,governance -masternode -txindex=1 -reindex -masternodeprivkey=91oHXFR2NVpUtBiJk37i8oBMChaQRbGjhnzWjN9KQ8LeAW7JBdN".format(PASTELD_DIR, EXTRACT_PATH_FULL),
"{}/pasteld -datadir={}/node8  -daemon -keypool=1 -discover=0 -rest -debug=masternode,mnpayments,governance -masternode -txindex=1 -reindex -masternodeprivkey=92MwGh67mKTcPPTCMpBA6tPkEE5AK3ydd87VPn8rNxtzCmYf9Yb".format(PASTELD_DIR, EXTRACT_PATH_FULL),
"{}/pasteld -datadir={}/node9  -daemon -keypool=1 -discover=0 -rest -debug=masternode,mnpayments,governance -masternode -txindex=1 -reindex -masternodeprivkey=92VSXXnFgArfsiQwuzxSAjSRuDkuirE1Vf7KvSX7JE51dExXtrc".format(PASTELD_DIR, EXTRACT_PATH_FULL),
"{}/pasteld -datadir={}/node10 -daemon -keypool=1 -discover=0 -rest -debug=masternode,mnpayments,governance -txindex=1 -reindex -gen".format(PASTELD_DIR, EXTRACT_PATH_FULL),
"{}/pasteld -datadir={}/node11 -daemon -keypool=1 -discover=0 -rest -debug=masternode,mnpayments,governance -txindex=1 -reindex -gen".format(PASTELD_DIR, EXTRACT_PATH_FULL),
"{}/pasteld -datadir={}/node12 -daemon -keypool=1 -discover=0 -rest -debug=masternode,mnpayments,governance -txindex=1 -reindex -gen".format(PASTELD_DIR, EXTRACT_PATH_FULL),
"{}/pasteld -datadir={}/node13 -daemon -keypool=1 -discover=0 -rest -debug=masternode,mnpayments,governance -txindex=1 -reindex -gen".format(PASTELD_DIR, EXTRACT_PATH_FULL),
"{}/pasteld -datadir={}/node14 -daemon -keypool=1 -discover=0 -rest -debug=masternode,mnpayments,governance -txindex=1 -reindex".format(PASTELD_DIR, EXTRACT_PATH_FULL),
"{}/pasteld -datadir={}/node15 -daemon -keypool=1 -discover=0 -rest -debug=masternode,mnpayments,governance -txindex=1 -reindex".format(PASTELD_DIR, EXTRACT_PATH_FULL),
"{}/pasteld -datadir={}/node16 -daemon -keypool=1 -discover=0 -rest -debug=masternode,mnpayments,governance -txindex=1 -reindex".format(PASTELD_DIR, EXTRACT_PATH_FULL),
"{}/pasteld -datadir={}/node17 -daemon -keypool=1 -discover=0 -rest -debug=masternode,mnpayments,governance -txindex=1 -reindex".format(PASTELD_DIR, EXTRACT_PATH_FULL),
"{}/pasteld -datadir={}/node18 -daemon -keypool=1 -discover=0 -rest -debug=masternode,mnpayments,governance -txindex=1 -reindex".format(PASTELD_DIR, EXTRACT_PATH_FULL)
]

def extractArchive():
    print ("Extracting nodes to: {}\n".format(EXTRACT_FOLDER))
            
    if ARCHIVE.endswith("tar.gz"):
        tar = tarfile.open(ARCHIVE, "r:gz")
        tar.extractall(path="./{}".format(EXTRACT_FOLDER))
        tar.close()
    else:
        print ("Archive type not supported!\n")

def start_pasteld():
    print ("Starting network")
    for command in CMD_LIST:
        print("Initiated command: {}".format(command))
        os.system(command)

if __name__ == '__main__':

    if( os.path.isdir("./"+EXTRACT_FOLDER) ):
        print ("Folder exists")
        if ( os.path.isdir( '{}/node0'.format("./"+EXTRACT_FOLDER) ) ):
            print ("No need to extract, nodes also exist\n")
        else:
            extractArchive()
    else:
        print ("Making data directory as : {}\n".format(ARCHIVE))
        os.mkdir("./"+EXTRACT_FOLDER)
        extractArchive()

    start_pasteld()
