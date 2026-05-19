# PicCapGenerator
**Amirkabir University of Technology | Fall 2024**  
**Course: Fundamentals of Cloud Computing**

## Overview
A microservices pipeline that accepts a user image, generates a text caption for it, and produces a similar image — delivering the result via email. Built in Go, deployed on cloud infrastructure.

## Architecture

```
User → Service 1 (Request Management)
           ↓ S3 Object Storage (image upload)
           ↓ MySQL (DBaaS)
           ↓ RabbitMQ (CloudAMQP)
                ↓
         Service 2 (Image Captioning)
           ↓ HuggingFace Captioning API
           ↓ MySQL (caption stored)
                ↓
         Service 3 (Image Generation)
           ↓ HuggingFace Text-to-Image API
           ↓ S3 Object Storage (generated image)
           ↓ MailerSend (email URL to user)
```

## Services

**Service 1 — Request Management**
- `POST /request` — accepts email + image, stores in S3, queues ID in RabbitMQ
- `GET /request/:id` — returns status (`pending` / `ready` / `done` / `failure`) or generated image URL

**Service 2 — Image Captioning**
- Consumes request ID from RabbitMQ
- Fetches image from S3, sends to HuggingFace captioning API
- Stores caption in DB, updates status to `ready`

**Service 3 — Image Generation**
- Polls DB for `ready` requests
- Sends caption to HuggingFace text-to-image API
- Uploads result to S3, emails public URL to user, updates status to `done`

## Cloud Services Used
| Purpose | Service |
|---------|---------|
| Cloud Host | Liara |
| Database (DBaaS) | MySQL on Liara |
| Object Storage | S3-compatible (ArvanCloud / Liara) |
| Message Queue | RabbitMQ via CloudAMQP |
| Image Captioning | HuggingFace Inference API |
| Image Generation | HuggingFace Text-to-Image API |
| Email | MailerSend |

## Setup
Copy `.env.example` to `.env` in each service and fill in credentials:
```bash
cp config/.env.example config/.env
```

Each service has its own `Dockerfile` for containerized deployment.

## DB Schema
```
Request(ID, Email, Status, ImageCaption, NewImageURL)
```
