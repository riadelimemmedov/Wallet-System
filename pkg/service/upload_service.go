package upload_service

/* The File Upload Flow

User uploads a file to the API
HandleFileUpload validates the file and saves it to temporary storage
An UploadJob is created with PENDING status and stored in the repository
The job is published to RabbitMQ exchange, which routes it to the upload queue
Worker(s) consume jobs from the queue and process them
When processing completes, the job status is updated to COMPLETED or FAILED
Users can query job status via GetUploadStatus

*/
