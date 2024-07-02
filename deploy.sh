#!/bin/bash

# Step 0: Build the HTML Files from MD
cd personal_site_builder && go run build_markdown.go ../entries ../personal_site_server/public
echo "Built HTML files from MD"

cd ../personal_site_server

# Step 1: Build the Go binary
echo "Building Go binary..."
GOOS=linux GOARCH=amd64 go build  -o personal_site ./server.go|| { echo "Go build failed"; exit 1; }

# Step 2: SCP the Go binary to the EC2 instance
echo "Copying Go binary to EC2 instance..."
scp -i /Users/jacksonstone/Desktop/Jackson\ Personal\ Site\ Key.pem personal_site ubuntu@3.19.146.227:/home/ubuntu/.temp/ || { echo "SCP failed"; exit 1; }

# Step 3: SSH into the EC2 instance and move the file
echo "Connecting to EC2 instance and moving the file..."
ssh -i /Users/jacksonstone/Desktop/Jackson\ Personal\ Site\ Key.pem ubuntu@3.19.146.227 << EOF
  mv ./.temp/personal_site . || { echo "Failed to move the file"; exit 1; }
  chmod +x personal_site || { echo "Failed to make the file executable"; exit 1; }
  echo "File moved successfully"
  sudo systemctl restart personal_site || { echo "Failed to restart"; exit 1; }
EOF

echo "Script completed successfully."
rm personal_site