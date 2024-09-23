# Homework 2

## Question 1

- What are packages in your implementation? What data structure do you use to transmit data and meta-data?

Packages is byte arrays with a size of 1024. These are what is used to transmit data and meta-data.

## Question 2

- Does your implementation use threads or processes? Why is it not realistic to use threads?

My implementation uses goroutines which is lightweight threads.

## Question 3

- In case the network changes the order in which messages are delivered, how would you handle message re-ordering?

Messages can be managed using a sequence number that defines the sequence of packages. We can then re-order messages
by this sequence number so it is continuous.

## Question 4

- In case messages can be delayed or lost, how does your implementation handle message loss?

My implementation does not handle message loss.

## Question 5

- Why is the 3-way handshake important?

The 3-way handshake is important beacause we want to ensure that the client and server er synced up and both
are ready to recive traffic. We also the server and client to agree on an initial squence number so we can prevent
data corruption.
