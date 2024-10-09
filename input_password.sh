#!/bin/sh
message="Please enter the password: "
while [ -z "$password" ]; do
    read -s -p "$message" password
    if [ -z "$password" ]; then
        echo -e "\033[2K"
        message="Password cannot be empty. Please enter the password: "
    fi
done
echo 
read -p "Please enter the helix relayer runner server address and port (default: 127.0.0.1:8081): " server
if [ -z "$server" ]; then
    server="127.0.0.1:8081"
fi

json_data=$(cat <<EOF
{
  "p": "$password"
}
EOF
)
unset password
echo "Sending password to $server..."
# print response code and response body
response=$(curl -s -w "\n%{http_code}" -X POST "http://$server/pass" \
     -H "Content-Type: application/json" \
     -d "$json_data")
# check response code is eq 200
if [ $(echo "$response" | tail -n1) -eq 200 ]; then
    echo "Send password success"
else
    echo "Send password failed, response: $response"
    exit 1
fi
