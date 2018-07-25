groupadd -g 1000 upspin
useradd -u 1000 -g 1000 -d / upspin
useradd -r -u 1000 -g upspin upspin

sudo -u upspin

./upspinserver -insecure -http :8080 -serverconfig /upspin-config/serverconfig