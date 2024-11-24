# go-expert-lab-stress-test
CLI em Go para realizar testes de carga em um serviço web.

Entrada de Parâmetros via CLI:

- -url: URL do serviço a ser testado.
- -requests: Número total de requests.
- -concurrency: Número de chamadas simultâneas.

### Buildar a imagem docker
```bash
    make build
```

### Testar

```bash
    docker container run go-http-bench "http://globo.com" -requests 10 -concurrency 2
```

## <a name="license"></a> License

Copyright (c) 2024 [Hugo Castro Costa]

[Hugo Castro Costa]: https://github.com/hgtpcastro
