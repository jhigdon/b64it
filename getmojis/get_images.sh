
curl https://slackmojis.com/ | grep img  > mojis.txt

#go get github.com/ericchiang/pup

cat mojis.txt | pup --color 'img json{}' > mojis.json





