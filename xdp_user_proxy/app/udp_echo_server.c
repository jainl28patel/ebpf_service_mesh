#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <unistd.h>
#include <arpa/inet.h>
#include <sys/types.h>
#include <sys/socket.h>

#define PORT 6000
#define BUFFER_SIZE 1024

int main() {
    // int sockfd;
    // struct sockaddr_in server_addr, client_addr;
    // socklen_t client_len = sizeof(client_addr);
    // char buffer[BUFFER_SIZE];
    // ssize_t recv_len;

    // // Create a UDP socket
    // if ((sockfd = socket(AF_INET, SOCK_DGRAM, 0)) < 0) {
    //     perror("socket creation failed");
    //     exit(EXIT_FAILURE);
    // }

    // // Prepare the sockaddr_in structure
    // memset(&server_addr, 0, sizeof(server_addr));
    // server_addr.sin_family = AF_INET;
    // server_addr.sin_addr.s_addr = INADDR_ANY; // Bind to all available interfaces
    // server_addr.sin_port = htons(PORT); // Use defined port

    // // Bind the socket to the port
    // if (bind(sockfd, (struct sockaddr *)&server_addr, sizeof(server_addr)) < 0) {
    //     perror("bind failed");
    //     close(sockfd);
    //     exit(EXIT_FAILURE);
    // }

    // printf("UDP server is running on port %d...\n", PORT);

    // // Main loop to receive and echo messages
    // while (1) {
    //     // Receive a message from a client
    //     recv_len = recvfrom(sockfd, buffer, BUFFER_SIZE, 0, (struct sockaddr *)&client_addr, &client_len);
    //     if (recv_len < 0) {
    //         perror("recvfrom failed");
    //         close(sockfd);
    //         exit(EXIT_FAILURE);
    //     }

    //     // Null-terminate the received data and print it
    //     buffer[recv_len] = '\0';
    //     printf("Received message: %s\n", buffer);

    //     // Send the same message back to the client
    //     if (sendto(sockfd, buffer, recv_len, 0, (struct sockaddr *)&client_addr, client_len) < 0) {
    //         perror("sendto failed");
    //         close(sockfd);
    //         exit(EXIT_FAILURE);
    //     }
    // }

    // // Close the socket (this will never be reached in this example)
    // close(sockfd);
    return 0;
}
