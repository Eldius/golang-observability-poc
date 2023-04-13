# configurando o ambiente num cluster k3s #

Decidi tentar configurar um cluster k8s numa Raspberry que tinha aqui de bobeira,
assim eu libero recursos na máquina local pra facilitar minha vida.


```shell
## colocar o conteudo abaixo no arquivo /boot/cmdline.txt
# cgroup_memory=1 cgroup_enable=memory
sudo vim /boot/cmdline.txt

## instalar o k3s
curl -sfL https://get.k3s.io | sh -

```

```shell
cat /etc/rancher/k3s/k3s.yaml

vim  ~/.kube/config
```

