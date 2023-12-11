# Public-5G NWDAF

This code was designed to run on the free5GC 5G core and its detailed description is located below.

## Tested Environment Configuration

As the execution environment has many components, the tested working versions of the main ones are listed down below:

- python: 3.8.10 (works with 3.7.x too)
- pip: 23.3.1
- tensorflow: 2.13.1
- keras: 2.13.1
- flask: 3.0.0
- go: 1.18.10
- free5GC: v3.3.0
- Ubuntu Server 20.04.3

**NOTE:** List updated on December, 2023

## Configuring the free5GC

Detailed instructions won't be added here as it isn't in the scope of this document, however as a general advice you should install the free5GC and then follow [these instructions](https://free5gc.org/guide/5-install-ueransim/#5-setting-free5gc-and-ueransim-parameters) configuring the IP address on AMF, SMF and UPF configuration files.

**TIP:** One should change the loopback IPs (127.0.0.x) to the one used by the LAN interface

## Install the prerequisites

**NOTE:** In this section the commands are supported on a BASH console

1. Install Python3 and pip
```
sudo apt install python3 python3-pip
pip install --upgrade pip
```

2. Install Python required packages
```
pip install -r requirements.txt
```

3. Install Flask
```
pip install flask
```

4. Install TensorFlow
```
pip install tensorflow
```

5. Add Go Lang support

**NOTE:** If Go version specified on [the first section](./README.md#tested-environment-configuration) is already installed, skip this step

**TIP:** Check Go version using the command `go version` on console

If another version of Go was previously installed, remove the existing version and install Go 1.18.10 using:
```
sudo rm -rf /usr/local/go
wget https://dl.google.com/go/go1.18.10.linux-amd64.tar.gz
sudo tar -C /usr/local -zxvf go1.18.10.linux-amd64.tar.gz
```
If not, install Go using the commands below:
```
wget https://dl.google.com/go/go1.18.10.linux-amd64.tar.gz
sudo tar -C /usr/local -zxvf go1.18.10.linux-amd64.tar.gz
mkdir -p ~/go/{bin,pkg,src}
# The following assumes that your shell is BASH:
echo 'export GOPATH=$HOME/go' >> ~/.bashrc
echo 'export GOROOT=/usr/local/go' >> ~/.bashrc
echo 'export PATH=$PATH:$GOPATH/bin:$GOROOT/bin' >> ~/.bashrc
echo 'export GO111MODULE=auto' >> ~/.bashrc
source ~/.bashrc
```

## Install the NWDAF module

### 1. Clone the repository inside the `free5GC` folder
```
cd ~/free5gc/
git clone -b mnc_Public-5G https://github.com/oliveiraleo/mnc_NWDAF.git
```

### 2. Copy the configuration file to the free5GC's `config` folder
```
cp mnc_NWDAF/nwdafcfg.yaml config/.
```

### 3. Load some required Go packages

First, change the working directory to the NWDAF's one

```
cd ~/free5gc/mnc_NWDAF/nwdaf
```
Then load the packages below
```
go mod download github.com/antonfisher/nested-logrus-formatter
go get nwdaf.com/service
go mod download github.com/free5gc/version
```

### 4. Install the NWDAF on free5GC

```
cd ~/free5gc/mnc_NWDAF/nwdaf
go build -o nwdaf nwdaf.go
cd ~/free5gc/
cp -r mnc_NWDAF/nwdaf .
```

## Running the Execution Environment

Execute the components in the following order:
1. 5G Core (free5GC)*
2. go module (nwdaf executable compiled on the previous section)
3. python module
4. temp_requester

\* NRF must be running

### To run NWDAF go module
```
cd ~/free5gc/nwdaf/
./nwdaf
# OR
go run nwdaf.go
```

Now the NWDAF should be running

### To run NWDAF python module
```
cd ~/free5gc/nwdaf/pythonmodule/
python main.py 
```

### To run temp_requester
```
cd ~/free5gc/nwdaf/Temp_Requester/
go run temp_requester.go
```

TODO: Format the text on the paragraph below or remove it

After that, you should select your number
If "1" is selected, MTLF (model training function) is executed.
Otherwise, "2" is selected, AnLF (analytics function) is executed.
Then, you can try to select a number which means the dataset number.
Now, we using the EMNIST dataset, which is in the python module.
In temp_requester, the image is not transmitted (using the json, the data number is transmitted).

## Configuring and Using UERANSIM

### 1. Install UERANSIM

On another machine (different from the free5GC one), run:
```
cd ~
git clone https://github.com/aligungr/UERANSIM
cd UERANSIM
sudo apt install make g++ libsctp-dev lksctp-tools iproute2
sudo snap install cmake --classic
make
```
### 2. Configure the Correct IP Addresses

```
cd ~/UERANSIM
nano config/free5gc-gnb.yaml
```

Find the section below:

```
ngapIp: 127.0.0.1   # gNB's local IP address for N2 Interface (Usually same with local IP)
gtpIp: 127.0.0.1    # gNB's local IP address for N3 Interface (Usually same with local IP)
```
And change `127.0.0.1` to the LAN IP from UERANSIM's machine

Now, on this line:

```
# List of AMF address information
amfConfigs:
  - address: 127.0.0.1
    port: 38412
```

Change `127.0.0.1` to the LAN IP from free5GC's machine

### 3. Add an UE to the core

Use the webconsole for that, the detailed instructions are [located here](https://free5gc.org/guide/5-install-ueransim/#4-use-webconsole-to-add-an-ue)

### 4. To run UERANSIM, use:

```
# gnb
build/nr-gnb -c config/free5gc-gnb.yaml
# ue
sudo build/nr-ue -c config/free5gc-ue.yaml
```

<!-- TODO: Finish merging the info below

### NWDAF Structure
NWDAF (Network Data Analytics Function) is consist of two part: 1) go module; 2) python module.

Go module can be run on "nwdaf.go" which located in "nwdaf" folder.

Python module can be run on "main.py" which located in "nwdaf/pythonmodule" folder.

### temp_requester Structure
temp_requester is the requester function which can be using on other NFs. 

If you want to call NWDAF from other NFs, the function in this requester can be used. -->
