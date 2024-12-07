🚀 OptiCuter Chat-Service
A comprehensive real-time chat service with support for video calls, review submissions, and video uploads — all underpinned by a robust, production-ready microservices architecture. Built with clean architecture principles for maintainability, scalability, and performance.

📚 Table of Contents
🌟 Features
🛠️ Tech Stack
📐 Architecture Overview
💾 Database Design
🚀 Installation & Usage
📈 Roadmap
🌟 Features
🔥 Real-Time Chat: Fast, efficient, and real-time chat functionality.
🎥 Video Calls: Jitsi-powered, high-quality, real-time video calling.
📹 Video Uploads: Chunked video uploads stored in MongoDB Atlas, with plans to integrate AWS S3 or Cloudinary for production.
✍️ Review Submissions: Collect and store feedback from users directly in the chat.
📡 gRPC Communication: Blazing-fast service-to-service communication.
⚙️ Concurrency Control: Leveraging Go routines, mutexes, and channels for smooth, bidirectional messaging.
🛠️ Tech Stack
Category	Technology	Description
Language	
High-performance back-end development.
Communication	
Seamless inter-service communication.
Video Calls	
Real-time video calls.
Storage	
Cloud storage for chat, video, and reviews.
Production	
Production-level video storage (coming soon).
Concurrency	Go Routines, Mutex, Channels	Advanced concurrency for bidirectional communication.
📐 Architecture Overview
                          +---------------------+   
                      |  API Gateway       |   
                      +----------+----------+   
                                 | 
        +------------------------+-----------------------+
        |                       gRPC                      |
        +------------+-------------------+---------------+
                     |                   |
          +----------+--------+  +-------+-----------+  
          |  Chat Service     |  |  Video Service   |  
          +-------------------+  +-----------------+  
                |                             |
         +------+---+               +---------+---+  
         | MongoDB  |               |   Jitsi Server  |
         +----------+               +-----------------+  

API Gateway: Unified access point for client requests.
Chat Service: Manages chats, reviews, and video uploads with MongoDB Atlas.
Video Service: Manages real-time video calls using Jitsi.
💾 Database Design
MongoDB Collections
Collection Name	Description
Users	Stores user profiles and authentication data.
Messages	Stores all chat messages, timestamps, and status.
Videos	Stores metadata for uploaded videos.
Reviews	Stores user-submitted reviews.
Sample MongoDB Document for Chat Messages:

json
Copy code
{
  "_id": "64c72f27d4f7f0a2b1234567",
  "senderId": "user123",
  "receiverId": "user456",
  "message": "Hello, how can I help you?",
  "timestamp": "2024-12-07T10:00:00Z",
  "status": "delivered"
}
🚀 Installation & Usage
Prerequisites
Go (v1.19 or later)
MongoDB Atlas Account
Docker (optional for containerized deployment)
Clone the Repository
git clone https://github.com/ratheeshkumar25/opti_cuter-chat.git
cd opti_cuter-chat
Set Up Environment Variables
Create a .env file with the following keys:

makefile

MONGODB_URI=<your_mongo_uri>
JITSI_API_URL=<your_jitsi_api_url>
S3_BUCKET_NAME=<your_s3_bucket_name> (if applicable)
Run the Service
bash

go run main.go
Test the API
Use tools like Postman or cURL to interact with the API Gateway.

📈 Roadmap
 ✅ Integrate AWS S3 / Cloudinary for production-level video uploads.
 🔄 WebSocket Support for real-time chat updates.
 📦 Dockerized Deployment for containerized CI/CD pipelines.
 📊 Performance Optimizations for chat, message, and video services.
 
