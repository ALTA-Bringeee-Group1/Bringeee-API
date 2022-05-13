<div id="top"></div>

<!-- PROJECT LOGO -->
<br/>
<div align="center">
<!--  mengarah ke repo  -->
  <a href="https://github.com/ALTA-Bringeee-Group1/Bringeee-API">
    <img src="images/logo.png" width="365" height="70">
  </a>

  <h3 align="center">Bringeee</h3>

  <p align="center">
    Final Project Capstone Program Immersive Alterra Academy
    <br />
    <a href="https://app.swaggerhub.com/apis-docs/wildanie12/Bringee-API/v1.2#/"><strong>Explore the docs Open API »</strong></a>
    <br />
    <br />
    <a href="https://www.codacy.com/gh/ALTA-Bringeee-Group1/Bringeee-API/dashboard?utm_source=github.com&amp;utm_medium=referral&amp;utm_content=ALTA-Bringeee-Group1/Bringeee-API&amp;utm_campaign=Badge_Grade">
      <img src="https://app.codacy.com/project/badge/Grade/3f9da093203f45a4b020bcedcb91196c">
    </a>
    <a href="https://www.codacy.com/gh/ALTA-Bringeee-Group1/Bringeee-API/dashboard?utm_source=github.com&utm_medium=referral&utm_content=ALTA-Bringeee-Group1/Bringeee-API&utm_campaign=Badge_Coverage">
      <img src="https://app.codacy.com/project/badge/Coverage/3f9da093203f45a4b020bcedcb91196c">
    </a>
  </p>
</div>

<br />

<!-- ABOUT THE PROJECT -->

## 💻 About The Project

Bringeee adalah sebuah aplikasi kargo untuk mengangkut barang dengan 1 truck. Harga untuk sebuah order dihitung berdasarkan jarak dari lokasi pengambilan order hingga ke tujuan pengiriman.

Berikut fitur yang terdapat dalam Bringeee :

<div>
      <details>
<summary>🙎 Customers</summary>
  
  <!---
  | Command | Description |
| --- | --- |
  --->
  
 Di Customer terdapat fitur untuk membuat Akun dan Login agar mendapat legalitas untuk mengakses berbagai fitur lain di aplikasi, 
 terdapat juga fitur Update untuk mengedit data yang berkaitan dengan customer, serta fitur delete berfungsi jika customer menginginkan hapus akun.
 
<div>
  
| Feature Customer | Endpoint | Param | JWT Token | Fungsi |
| --- | --- | --- | --- | --- |
| POST | /api/customers  | - | NO | Melakukan proses registrasi customer |
| POST | /api/auth | - | NO | Melakukan proses login customer |
| GET | /api/auth/me | - | YES | Mendapatkan informasi customer yang sedang login |
| PUT | /api/customers | - | YES | Melakukan update informasi customer yang sedang login | 
| DEL | /api/customers | - | YES | Menghapus customer yang sedang login |
| POST | /api/customers/orders | - | YES | Membuat sebuah order |
| POST | /api/customers/orders/estimate | - | YES | Melihat perkiraan harga sebuah orderan |
| GET | /api/customers/orders | status order | YES | Mendapatkan semua order berdasarkan status order customer |
| GET | /api/customers/orders/{orderID} | orderID | YES | Mendapatkan detail sebuah order customer |
| GET | /api/customers/orders/{orderID}/histories | orderID | YES | Mendapatkan timeline dari sebuah order |
| POST | /api/customers/orders/{orderID}/confirm | orderID | YES | Mengkonfirmasi sebuah order jika terjadi penyesuaian harga dari admin |
| POST | /api/customers/orders/{orderID}/cancel | orderID | YES | Membatalkan sebuah order |
| POST | /api/customers/orders/{orderID}/payment | orderID | YES | Memilih jenis pembayaran yang akan digunakan |
| POST | /api/customers/orders/{orderID}/payment/cancel | orderID | YES | Membatalkan jenis pembayaran |
| GET | /api/customers/orders/{orderID}/payment | orderID | YES | Mendapatkan jenis pembayaran yang digunakan oleh customer |

</details>

<details>
<summary>🚚 Driver</summary>
  
  <!---
  | Command | Description |
| --- | --- |
  --->
  
 Di Driver terdapat fitur untuk membuat Akun dan Login agar mendapat legalitas untuk mengakses berbagai fitur lain di aplikasi, 
 terdapat juga fitur Update untuk mengedit data yang berkaitan dengan driver, serta fitu - fitur lainnya.
 
<div>
  
| Feature Driver | Endpoint | Param | JWT Token | Fungsi |
| --- | --- | --- | --- | --- |
| POST | /api/drivers  | - | NO | Melakukan proses registrasi driver |
| POST | /api/auth | - | NO | Melakukan proses login driver |
| GET | /api/auth/me | - | YES | Mendapatkan informasi driver yang sedang login |
| PUT | /api/drivers | - | YES | Melakukan update informasi yang tidak credential driver yang sedang login | 
| GET | /api/drivers/orders | - | YES | Mendapatkan semua order berdasarkan tipe truck driver |
| GET | /api/drivers/current_order | - | YES | Mendapatkan order yang sedang diantar oleh driver |
| GET | /api/drivers/history_orders | - | YES | Mendapatkan order yang telah diantar oleh driver |
| GET | /api/orders/{orderID} | orderID | YES | Mendapatkan detail sebuah order |
| POST | /api/drivers/orders/{orderID}/take_order | orderID | YES | Mengambil sebuah orderan untuk diantarkan ke tujuan |
| POST | /api/drivers/orders/{orderID}/finish_order | orderID | YES | Menyelesaikan sebuah orderan dengan mengupload foto pada saat diterima customer |

</details> 
</details>

<details>
<summary>🖥️ Admin</summary>
  
  <!---
  | Command | Description |
| --- | --- |
  --->
  
 Di Admin terdapat fitur untuk melakukan manajemen customer, driver, order dan fitur statistik serta laporan order perbulannya.
 
<div>
  
| Feature Admin | Endpoint | Param | JWT Token | Fungsi |
| --- | --- | --- | --- | --- |
| POST | /api/auth | - | NO | Melakukan proses login admin |
| GET | /api/auth/me | - | YES | Mendapatkan informasi admin yang sedang login |
| GET | /api/customers | (optional) | YES | Mendapatkan list customer |
| GET | /api/drivers | (optional) | YES | Mendapatkan list driver |
| GET | /api/orders | (optional) | YES | Mendapatkan list order |
| GET | /api/orders/{orderID}/histories | orderID | YES | Mendapatkan timeline sebuah order |
| POST | /api/orders/{orderID}/confirm | orderID | YES | Mengkonfirmasi sebuah order jika tidak ada penyesuaian harga |
| POST | /api/orders/{orderID}/cancel | orderID | YES | Membatalkan sebuah order |
| GET | /api/orders/{orderID} | orderID | YES | Mendapatkan detail sebuah order |
| PATCH| /api/orders/{orderID} | orderID | YES | Melakukan penyesuaian harga pada sebuah order | 
| POST | /api/drivers/{driverID}/confirm | driverID | YES | Mengverifikasi akun driver |
| GET | /api/drivers/{driverID} | driverID | YES | Mendapatkan detail profile driver |
| PUT | /api/drivers/{driverID} | driverID | YES | Melakukan update informasi yang credential pada akun driver | 
| DEL | /api/drivers/{driverID} | driverID | YES | Menghapus akun driver | 
| GET | /api/customers/{customerID} | customerID | YES | Mendapatkan detail profile customer |
| DEL | /api/customers/{customerID} | customerID | YES | Menghapus akun customer | 
| GET | /api/stats/aggregates/orders_count | (optional) | YES | Mendapatkan jumlah semua order |
| GET | /api/stats/aggregates/drivers_count | (optional) | YES | Mendapatkan jumlah semua driver |
| GET | /api/stats/aggregates/truck_types_count | - | YES | Mendapatkan jumlah semua tipe truck |
| GET | /api/stats/aggregates/customers_count | - | YES | Mendapatkan jumlah semua customer |
| GET | /api/stats/orders/{day} | day | YES | Mendapatkan jumlah order berdasarkan periode hari yang di inginkan |
| POST | /api/export/orders | - | YES | Membuat file excel laporan order perbulan |

</details>

<p align="right">(<a href="#top">back to top</a>)</p>

### Built With

### 🛠 &nbsp;Build App & Database

![JSON](https://img.shields.io/badge/-JSON-05122A?style=flat&logo=json&logoColor=000000)&nbsp;
![GitHub](https://img.shields.io/badge/-GitHub-05122A?style=flat&logo=github)&nbsp;
![Visual Studio Code](https://img.shields.io/badge/-Visual%20Studio%20Code-05122A?style=flat&logo=visual-studio-code&logoColor=007ACC)&nbsp;
![MySQL](https://img.shields.io/badge/-MySQL-05122A?style=flat&logo=mysql&logoColor=4479A1)&nbsp;
![Golang](https://img.shields.io/badge/-Golang-05122A?style=flat&logo=go&logoColor=4479A1)&nbsp;
![Echo](https://img.shields.io/badge/-Echo-05122A?style=flat&logo=go)&nbsp;
![Gorm](https://img.shields.io/badge/-Gorm-05122A?style=flat&logo=go)&nbsp;
![AWS](https://img.shields.io/badge/-AWS-05122A?style=flat&logo=amazon)&nbsp;
![Insomnia](https://img.shields.io/badge/-Insomnia-05122A?style=flat&logo=insomnia)&nbsp;
![Docker](https://img.shields.io/badge/-Docker-05122A?style=flat&logo=docker)&nbsp;
![Ubuntu](https://img.shields.io/badge/-Ubuntu-05122A?style=flat&logo=ubuntu)&nbsp;
![Midtrans](https://img.shields.io/badge/-Midtrans-05122A?style=flat&logo=midtrans)&nbsp;
![Codacy](https://img.shields.io/badge/-Codacy-05122A?style=flat&logo=codacy)&nbsp;
![CloudFlare](https://img.shields.io/badge/-CloudFlare-05122A?style=flat&logo=cloudflare)&nbsp;
![Google Maps Platform](https://img.shields.io/badge/-Google_Maps_Platform-05122A?style=flat&logo=google)&nbsp;
![JWT](https://img.shields.io/badge/-JWT-05122A?style=flat&logo=jwt)&nbsp;
![Swagger](https://img.shields.io/badge/-Swagger-05122A?style=flat&logo=swagger)&nbsp;
![Lucid Chart](https://img.shields.io/badge/-Lucid_Chart-05122A?style=flat&logo=lucidchart)&nbsp;

<p align="right">(<a href="#top">back to top</a>)</p>

## 🗃️ ERD

<img src="images/erd.png">
<p align="right">(<a href="#top">back to top</a>)</p>

## 📑 Use Case Diagram

<img src="images/UCD.png">

<p align="right">(<a href="#top">back to top</a>)</p>

<!-- CONTACT -->

## 📱 Contact

Project Repository Link : [https://github.com/ALTA-Bringeee-Group1/Bringeee-API](https://github.com/ALTA-Bringeee-Group1/Bringeee-API)<br/>
Open API Documentation : [https://app.swaggerhub.com/apis-docs/wildanie12/Bringee-API/v1.2#/](https://app.swaggerhub.com/apis-docs/wildanie12/Bringee-API/v1.2#/)&nbsp;

<!-- :heart: -->
<!-- CONTRIBUTOR -->

Contributor :
<br>
[![Gmail: M. Badar Wildani](https://img.shields.io/badge/-badar.wildanie@gmail.com-maroon?style=flat&logo=gmail)](https://mail.google.com/mail/u/0/#inbox?compose=CllgCJqTfrDgzWPFFgSKDLmBlPGRmCRXMQVTgqZDWJrxHDMJkSBGGCGnnGJhRKjrbzjJmFqnZFg)
[![GitHub M. Badar Wildani](https://img.shields.io/github/followers/wildanie12?label=follow&style=social)](https://github.com/wildanie12)

[![Gmail: Muh. Nasrul](https://img.shields.io/badge/-nasrulmuhammad748@gmail.com-maroon?style=flat&logo=gmail)](https://mail.google.com/mail/u/0/#inbox?compose=CllgCJqTfrDgzWPFFgSKDLmBlPGRmCRXMQVTgqZDWJrxHDMJkSBGGCGnnGJhRKjrbzjJmFqnZFg)
[![GitHub Muh. Nasrul](https://img.shields.io/github/followers/mnasruls?label=follow&style=social)](https://github.com/mnasruls)

<p align="right">(<a href="#top">back to top</a>)</p>
<h3>
<p align="center">:copyright: 2022 | Built with :heart: from us</p>
</h3>
<!-- end -->
