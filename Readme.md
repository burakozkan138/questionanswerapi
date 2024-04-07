**Question Answer API**

Bu proje, bir soru-cevap API'si geliştirmeyi amaçlamaktadır. API, kullanıcılara soru sorma ve bu sorulara cevap alma imkanı sunar. Ayrıca, API kullanıcılarının sorularını kaydetme, güncelleme ve silme işlemlerini gerçekleştirebilir.

## Kurulum
Aşağıdaki adımları izleyerek projeyi yerel ortamınızda çalıştırabilirsiniz.

### 1. Adım
GitHub deposunu yerel ortamınıza klonlayın:
```bash
git clone https://github.com/burakozkan138/questionanswerapi.git
```

### 2. Adım
Projenin root dizinine gidin.
```bash
cd questionanswerapi
```

### 3. Adım
Gerekli config dosyalarını ilk olarak örnek dosyadan kopyalayın. Daha sonra, kopyaladığınız dosyaları açarak gerekli düzenlemeleri yapın.
```bash
copy config\.env.example config\.env && copy config\.env.example config\.env.test
```

### 4. Adım
Config dosyasında gerekli ayarlamaları yaptıktan sonra docker-compose.yml üzerinde database environment alanı kontrol edildikten sonra aşağıdaki komut çalıştırılarak proje docker üzerinde ayağa kaldırılır.
```bash
docker-compose up -d --build
```

## Kullanım
Proje başarı ile ayağa kalktı ise swagger arayüzüne erişmek için aşağıdaki URL'yi ziyaret edin:
```bash
http://localhost:8080/swagger
```

### Şuan için swagger üzerinde authentication alanı gereken requestlerde bazı hatalar mevcut ancak postman ile oluşturduğum workspace üzerinde herhangibir hata görmedim.
```bash
https://www.postman.com/cinemabookingsystem/workspace/questionanswer
```
