# 📃 FreeBox Scanner Documentation

## Table of Contents
- [📃 FreeBox Scanner Documentation](#-freebox-scanner-documentation)
	- [Table of Contents](#table-of-contents)
	- [Introduction](#introduction)
	- [Usage](#usage)
	- [Requirements](#requirements)
	- [Installation](#installation)
	- [Troubleshooting](#troubleshooting)
	- [Responsible Use](#responsible-use)
	- [Free SAS CIDR range list](#free-sas-cidr-range-list)
	- [See PoC on Youtube](#see-poc-on-youtube)

## Introduction
The Freebox Port Scanner and Exploiter is a powerful tool designed to scan for Freebox routers on a network and identify those with open ports susceptible to exploitation. The Freebox is a popular router model used in many homes and offices, and while it provides essential networking features, it may have certain vulnerabilities that can be exploited if not properly secured.

This tool serves as a security assessment utility to help individuals and organizations identify potential weaknesses in their Freebox routers. By scanning for open ports, the tool can pinpoint services and protocols that might be unintentionally exposed to the public internet, increasing the risk of unauthorized access and cyberattacks.

It is essential to note that the primary purpose of this tool is to raise awareness and promote responsible security practices. Using this tool against systems you do not own or without explicit permission is illegal and unethical. Always ensure you have proper authorization before performing any security testing.

## Usage
To use the scanner, follow these steps:
1. [Download](#installation) and install the scanner on the source system.
2. Run the scan using free CIDR.

## Requirements
Before installing and using the Binary Dropper, ensure that the following requirements are met:
- Operating System: Windows 10/11, Linux
- Golang 1.20

## Installation
1. Clone or download the Binary Dropper repository from [Here](https://github.com/0xF7A4C6/freebox-scan/tree/main).
2. Navigate to the downloaded directory.
3. Configure the add cidr range into `assets/asn.txt` .
4. Install make and Build `Makefile` (`make scan, make brute, make build`).
5. Add proxies into `assets/asn.txt`.

## Troubleshooting
If you encounter any issues while using the Binary Dropper, try the following troubleshooting steps:

1. Update golang
2. Make sur to set working proxies, asn

If the problem persists, feel free to open an issue on the GitHub repository for support.

## Responsible Use

It is essential to use the Freebox Port Scanner and Exploiter responsibly and ethically. Unauthorized use against systems you do not own or have explicit permission to scan is illegal and may result in severe legal consequences. Always ensure you have proper authorization and consent before using this tool on any network.

Remember, the goal is not to cause harm but to raise awareness about potential security risks. If you discover vulnerabilities or open ports on your own Freebox router, take immediate action to secure your network. If you encounter vulnerabilities on other systems, report them responsibly to the affected parties or follow proper disclosure procedures.

Proceed with caution, act responsibly, and prioritize the security and privacy of others while using this tool.

## Free SAS CIDR range list

```txt
88.176.0.0/12	Free SAS	France cidr France,Limeil-Brevannes
91.160.0.0/12	Free SAS	France cidr France,Saint-Médard-en-Jalles
82.248.0.0/13	Free SAS	France cidr France,Courbevoie
88.176.0.0/13	Free SAS	France cidr France,Limeil-Brevannes
78.224.0.0/13	Free SAS	France cidr France,Geneuille
91.160.0.0/14	Free SAS	France cidr France,Saint-Médard-en-Jalles
88.168.0.0/14	Free SAS	France cidr France,Tournon-sur-Rhône
83.152.0.0/14	Free SAS	France cidr France,Troyes
88.124.0.0/14	Free SAS	France cidr France,Labarthe-sur-Leze
88.184.0.0/14	Free SAS	France cidr France,Montrabe
88.160.0.0/14	Free SAS	France cidr France,Paris
83.156.0.0/14	Free SAS	France cidr France,Blagnac
88.172.0.0/15	Free SAS	France cidr France,Les Arcs
78.232.0.0/15	Free SAS	France cidr France,Valence
88.188.0.0/15	Free SAS	France cidr France,Roubaix
88.174.0.0/15	Free SAS	France cidr France,Gommerville
88.166.0.0/15	Free SAS	France cidr France,Béziers
88.164.0.0/16	Free SAS	France cidr France,Dreux
78.254.0.0/16	Free SAS	France cidr France
88.120.128.0/17	Free SAS	France cidr France,Mandelieu-la-Napoule
```

## See PoC on Youtube

[![Preview](https://img.youtube.com/vi/3SKR2KtgTGU/0.jpg)](https://youtu.be/3SKR2KtgTGU)
