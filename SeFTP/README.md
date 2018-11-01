# SeFTP | [![pipeline status](https://gitlab.com/clover/SeFTP/badges/Development/pipeline.svg)](https://gitlab.com/clover/SeFTP/commits/Development) | [![Go Report Card](https://goreportcard.com/badge/gitlab.com/clover/SeFTP)](https://goreportcard.com/report/gitlab.com/clover/SeFTP)

SeFTP is a fast, secure and reliable file transfer protocol written in Golang.

SeFTP is designed to transfer file securely and fastly. 
It takes advantage of Authenticated Encryption with Associated Data (AEAD) algorithm, UDP stream transfer, and TCP(possibly include obfuscation algorithm) to achieve security, speed, and stability simultaneously. 
SeFTP constructs on the reliability of TCP and the high speed of UDP(Even with TCP congestion algorithm, TCP is still slower in most cases). 
The main authentication part is done during TCP as it requires reliability and does not require much speed. 
The file transfer part can be done in either UDP or TCP to take advantage of different protocols.

SeFTP is under HEAVY development, so bugs MUST vary. But, anyone is welcomed to contribute by forking this repo.