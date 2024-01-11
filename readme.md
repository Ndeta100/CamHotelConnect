# Hotel Reservation API

## Overview

This project is a backend JSON API for a hotel reservation system. It's built using Go (Golang) and utilizes MongoDB for data persistence. The API allows for room booking at a hotel, with features to support both user and administrative functionalities.

## Features

### For Users

- **Room Booking**: Users can browse available rooms and make reservations.
- **Account Management**: Users can create accounts, log in, and manage their bookings.

### For Admins

- **Booking Management**: Admins can view and manage all room bookings.
- **Room Management**: Admins have the capabilities to add, remove, or modify room details.
- **User Management**: Admins can manage user accounts and access rights.

### Security

- **Authentication**: The API supports secure user authentication.
- **Authorization**: Different authorization levels are implemented, ensuring users have access only to appropriate functionalities.

### Additional Utilities

- **Database Seeding**: Scripts are provided to seed the database with initial data.
- **Database Migrations**: The system supports database migrations for smooth transitions and upgrades.

## Technologies Used

- **Programming Language **: [Go (Fiber)](https://gofiber.io/)
- **Database**: [MongoDB](https://www.mongodb.com/docs/drivers/go/current/)
- **Authentication**: [JWT](https://jwt.io/) / [OAuth](https://oauth.net/2/) (TBD)

## Getting Started

### Prerequisites

- [Go (latest version)](https://golang.org/dl/)
- [MongoDB](https://www.mongodb.com/try/download/community)

### Installation

1. **Clone the repository**:
   ```bash
   git clone https://github.com/Ndeta100/CamHotelConnect.git
   
### project env variables
```
   HTTP_LISTEN_ADDRESS=:3000
   JWT_SECRET=somthing_supersecret_No_bodyknows
   MONGO_DB_URL=mongodb://localhost:27017
   MONGO_DB_NAME=hotel-reservation
   MONGO_DB_URL_TEST=mongodb://localhost:27017
```
