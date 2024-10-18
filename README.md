# PI Installation steps:

### Brief overview just for ref

1. Install Raspberry PI OS Lite
    - Enable SSH

2. Update / Upgrade the OS

3. Install the Mosquitto MQTT Broker
    - Test the Broker

4. Set Up Zigbee2MQTT
    - If using serial over usb setup serial before installing ```sudo raspi-config``` and update ```Interface Options > Serial > No > Yes > sudo reboot```
    - Download node