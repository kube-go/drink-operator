# drink-operator

A simple kubernetes operator

## Run controller in one terminal

```go
go run main.go or make run
```

## Create custom resources for tea in another terminal

```go
kubectl apply -f config/samples/hotdrink_v1alpha1_tea.yaml
```

```go
░▒▓    ~/github.com/kube-go/drink-operator  on   main ⇡2 !7  kubectl apply -f config/samples/hotdrink_v1alpha1_tea.yaml                ✔  3.9.5   1.16.4   2.7.3   at kind-pai/crd ⎈  at 12:51:26  ▓▒░
tea.hotdrink.slashpai.github.io/black-tea created
tea.hotdrink.slashpai.github.io/milk-tea created
tea.hotdrink.slashpai.github.io/milk-tea-no-sugar created

░▒▓    ~/gi/k/drink-operator  on   main ⇡2 !7  k get tea                                                                               ✔  3.9.5   1.16.4   2.7.3   at kind-pai/crd ⎈  at 12:53:06  ▓▒░
NAME                AGE
black-tea           82s
milk-tea            82s
milk-tea-no-sugar   82s

░▒▓    ~/gi/k/drink-operator  on   main ⇡2 !7  k get pods                                                                              ✔  3.9.5   1.16.4   2.7.3   at kind-pai/crd ⎈  at 12:54:28  ▓▒░
NAME                    READY   STATUS      RESTARTS   AGE
black-tea-pod           0/1     Completed   0          85s
milk-tea-no-sugar-pod   0/1     Completed   0          84s
milk-tea-pod            0/1     Completed   0          84s

░▒▓    ~/gi/k/drink-operator  on   main ⇡2 !7  
```

```
░▒▓    ~/gi/k/drink-operator  on   main ⇡2 !7  k logs black-tea-pod                                                                    ✔  3.9.5   1.16.4   2.7.3   at kind-pai/crd ⎈  at 12:54:31  ▓▒░
make black tea without sugar

░▒▓    ~/gi/k/drink-operator  on   main ⇡2 !7  k logs milk-tea-no-sugar-pod                                                            ✔  3.9.5   1.16.4   2.7.3   at kind-pai/crd ⎈  at 13:02:17  ▓▒░
make milk tea without sugar

░▒▓    ~/gi/k/drink-operator  on   main ⇡2 !7  k logs milk-tea-pod                                                                     ✔  3.9.5   1.16.4   2.7.3   at kind-pai/crd ⎈  at 13:02:23  ▓▒░
make milk tea with sugar

░▒▓    ~/gi/k/drink-operator  on   main ⇡2 !7 
```
