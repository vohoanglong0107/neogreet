# Neogreet

## SELinux

```bash
cat audit.log | audit2allow -M neogreet
sudo semodule -X 300 -i neogreet.pp
```
