## a boot OnDefineDomain sidecar

### domain.os.boot

```diff
apiVersion: kubevirt.io/v1
kind: VirtualMachine
spec:
  template:
    metadata:
      annotations:
+       hooks.kubevirt.io/hookSidecars: '[{"args": ["--version", "v1alpha3"],"image": "ghcr.io/go-bai/kubevirt-boot-sidecar:v1.2.0"}]'
+       os.vm.kubevirt.io/boot: '{"boot":[{"dev":"hd"},{"dev":"cdrom"}]}'
```

### domain.os.bootmenu

```diff
apiVersion: kubevirt.io/v1
kind: VirtualMachine
spec:
  template:
    metadata:
      annotations:
+       hooks.kubevirt.io/hookSidecars: '[{"args": ["--version", "v1alpha3"],"image": "ghcr.io/go-bai/kubevirt-boot-sidecar:v1.2.0"}]'
+       os.vm.kubevirt.io/boot: '{"bootmenu":{"enable": "yes", "timeout": "10000"}}'
```
