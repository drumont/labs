

```bash
wget -O - https://apt.releases.hashicorp.com/gpg | sudo gpg --dearmor -o /usr/share/keyrings/hashicorp-archive-keyring.gpg
echo "deb [arch=$(dpkg --print-architecture) signed-by=/usr/share/keyrings/hashicorp-archive-keyring.gpg] https://apt.releases.hashicorp.com $(lsb_release -cs) main" | sudo tee /etc/apt/sources.list.d/hashicorp.list
sudo apt update && sudo apt install vault

systemctl enable vault.service
systemctl start vault.service

```



### Modify /etc/vault.d/vault.hcl
```bash
ui = true
cluster_addr  = "http://127.0.0.1:8201"
api_addr      = "https://127.0.0.1:8200"
disable_mlock = true

#mlock = true
#disable_mlock = true

storage "file" {
  path = "/opt/vault/data"
}

#storage "consul" {
#  address = "127.0.0.1:8500"
#  path    = "vault"
#}

# HTTP listener
#listener "tcp" {
#  address = "127.0.0.1:8200"
#  tls_disable = 1
#}

# HTTPS listener
listener "tcp" {
  address       = "0.0.0.0:8200"
  cluster_addr  = "0.0.0.0:8201"
  tls_cert_file = "/opt/vault/tls/tls.crt"
  tls_key_file  = "/opt/vault/tls/tls.key"
}

# Enterprise license_path
# This will be required for enterprise as of v1.8
#license_path = "/etc/vault.d/vault.hclic"

# Example AWS KMS auto unseal
#seal "awskms" {
#  region = "us-east-1"
#  kms_key_id = "REPLACE-ME"
#}

# Example HSM auto unseal
#seal "pkcs11" {
#  lib            = "/usr/vault/lib/libCryptoki2_64.so"
#  slot           = "0"
#  pin            = "AAAA-BBBB-CCCC-DDDD"
#  key_label      = "vault-hsm-key"
#  hmac_key_label = "vault-hsm-hmac-key"
#}
```

### Verify configuration
```bash
vault server -dev -config /etc/vault.d/vault.hcl
```


vault write pki/root/generate/internal \
    common_name=intra.nathos.dev \
    ttl=8760h


vault write pki/config/urls \
    issuing_certificates="http://192.168.1.32:8200/v1/pki/ca" \
    crl_distribution_points="http://192.168.1.32:8200/v1/pki/crl"

vault write pki/roles/intra-nathos-pki-role \
    allowed_domains=intra.nathos.dev \
    allow_subdomains=true \
    max_ttl=72h

vault policy write pki - <<EOF
path "pki*"                        { capabilities = ["read", "list"] }
path "pki/sign/intra-nathos-pki-role"    { capabilities = ["create", "update"] }
path "pki/issue/intra-nathos-pki-role"   { capabilities = ["create"] }
EOF


