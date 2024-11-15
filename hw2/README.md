# Homework 2

## Question 1

- What are packages in your implementation? What data structure do you use to transmit data and meta-data?

A package is a custom structure that has both metadata fields such as flags and the data it holds.

## Question 2

- Does your implementation use threads or processes? Why is it not realistic to use threads?

My implementation uses threads which is lightweight threads.

## Question 3

- In case the network changes the order in which messages are delivered, how would you handle message re-ordering?

Messages can be managed using a sequence number that defines the sequence of packages. We can then re-order messages by this sequence number so it is continuous or detect data loss if it isn't.

## Question 4

- In case messages can be delayed or lost, how does your implementation handle message loss?

My implementation resets the connection sending a package with the RST flag.

## Question 5

- Why is the 3-way handshake important?

The 3-way handshake is important because we want to ensure that the client and server are synced up and both are ready to receive traffic. We also want to make sure the server and client to agree on an initial sequence number so we can prevent data corruption.
