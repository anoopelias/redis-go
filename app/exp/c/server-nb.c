#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <sys/socket.h>
#include <arpa/inet.h>
#include <unistd.h>
#include <fcntl.h>
#include <errno.h>

int main(int argc, char *argv[]) {
    int socket_desc, client_sock, c, read_size;
    struct sockaddr_in server, client;
    char client_message[2000];

    // Create socket
    socket_desc = socket(AF_INET, SOCK_STREAM, 0);
    if (socket_desc == -1) {
        printf("Could not create socket");
    }
    puts("Socket created");

    // Prepare the sockaddr_in structure
    server.sin_family = AF_INET;
    server.sin_addr.s_addr = INADDR_ANY;
    server.sin_port = htons(8888);

    // Bind
    if (bind(socket_desc, (struct sockaddr *)&server, sizeof(server)) < 0) {
        perror("Bind failed");
        return 1;
    }
    puts("Bind done");

    // Listen
    if (listen(socket_desc, 3) < 0) {
        perror("Listen failed");
        return 1;
    }
    puts("Waiting for incoming connections...");

    // Set the socket to non-blocking mode
    int flags = fcntl(socket_desc, F_GETFL, 0);
    fcntl(socket_desc, F_SETFL, flags | O_NONBLOCK);

    while (1) {
        c = sizeof(struct sockaddr_in);
        client_sock = accept(socket_desc, (struct sockaddr *)&client, (socklen_t *)&c);
        if (client_sock < 0) {
            if (errno == EWOULDBLOCK || errno == EAGAIN) {
                // No incoming connections at this time, wait a bit
                usleep(1000);
            }
            else {
                perror("Accept failed");
                return 1;
            }
        }
        else {
            puts("Connection accepted");

            // Set the client socket to non-blocking mode
            int flags = fcntl(client_sock, F_GETFL, 0);
            fcntl(client_sock, F_SETFL, flags | O_NONBLOCK);

            // Receive messages from the client
            while (1) {
                read_size = recv(client_sock, client_message, 2000, 0);
                if (read_size < 0) {
                    if (errno == EWOULDBLOCK || errno == EAGAIN) {
                        // No more messages at this time, wait a bit
                        usleep(1000);
                    }
                    else {
                        perror("Receive failed");
                        break;
                    }
                }
                else if (read_size == 0) {
                    puts("Client disconnected");
                    break;
                }
                else {
                    // Print the received message
                    printf("Client message: %s", client_message);
                }
            }

            // Close the client socket
            close(client_sock);
        }
    }

    return 0;
}
