# CCS-E4003 - Special Assignment in Computer Science

# Customizing Kubernetes Container Scheduler on IoT Nodes

The aim of this special assignment is to experiment with custom Kubernetes container scheduler for IoT environment 

## Team members:
  - Rajagopalan Ranganathan
  - Sunil Kumar Mohanty

## Project Description
### Overview

The aim of the project is to experiment, understand and build a custom Kubernetes container scheduler for Internet of Things (IoT) nodes. IoT environment is highly dynamic and heterogenous. The nodes could vary drastically with respect to processing power, hardware and other factors.
IoT nodes could be grouped based on their Geo location, for example a set of nodes could be present in "Home" and a few others in the "farm" and henceforth. Using a custom container scheduler to instruct which PODs would be deployed in which nodes based on the above said parameters eases Node and POD management. 



### High Level Architecture

The high level architecture of the solution is depicted in the Figure below.

![Kubernetes http://k8s.info/cs.html](readme-res/k-arch.png "Kubernetes (http://k8s.info/cs.html)")
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;Kubernetes (http://k8s.info/cs.html)


### Custom Scheduler

The custom scheduler created for this purpose, is written in "GoLang". The scheduler keeps running, scanning for new Nodes and Pods and assigning the Pods to the nodes based on a custom logic. The logic used, is to associate the POD to a corresponding Node based on the location category (Ex: Home, Office, Farm) and the network requirements (Ex: Ethernet, Wifi), and randomly chosing a Node if the requirements are not met absolutely.


### Environment

Kubernetes --

Nodes  -- Ubuntu 16.04

GoLang - 1.8

Dcoker --

