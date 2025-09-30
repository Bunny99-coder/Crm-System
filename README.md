# Real Estate CRM

  A full-stack CRM application designed for real estate professionals to manage properties, leads, deals, and client interactions
  efficiently. The backend is built with Go, and the frontend is a modern web application using Next.js and TypeScript.

  ## ✨ Features

     Dashboard:* At-a-glance overview of key metrics, tasks, and deals.
     Contact Management:* A centralized database for all clients and leads.
     Deal Tracking:* Visualize and manage the sales pipeline from lead to closing.
     Property Listings:* Manage property information, status, and associated deals.
     Task Management:* Create and assign tasks to stay on top of daily activities.
     Reporting:* Generate reports to analyze sales performance and trends.
     User Authentication:* Secure login and registration for different user roles (e.g., Sales Agent, Reception).

  ## 🛠️ Tech Stack

     Backend:*
      *   Go
      *   PostgreSQL
      *   chi router (likely, based on common Go API patterns)
      *   sqlx for database interaction (likely)
     Frontend:*
      *   Next.js (React)
      *   TypeScript
      *   Tailwind CSS
      *   Shadcn/ui for components
     Package Manager:* pnpm

  ## 🚀 Getting Started

  Follow these instructions to set up and run the project on your local machine.

  ### Prerequisites

  *   Go (https://go.dev/doc/install) (version 1.21 or later)
  *   Node.js (https://nodejs.org/en) (version 20.x or later)
  *   pnpm (https://pnpm.io/installation)
  *   PostgreSQL (https://www.postgresql.org/download/)
  *   A running PostgreSQL instance.

  ### 1. Backend Setup

  The backend server is written in Go and connects to a PostgreSQL database.

  `bash
  # 1. Navigate to the backend directory
  cd backend

  # 2. Install Go dependencies
  go mod tidy

  3. Set up environment variables
  server:
    port: 8080
  `

  `bash
  # 4. Run database migrations
  # You will need a migration tool like 'golang-migrate/migrate'.
  # If not installed, you can install it via:
  # go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest


  migrate -database 'postgres://your_db_user:your_db_password@localhost:5432/real_estate_crm?sslmode=disable' -path db/migrations up


  5. Run the server
  3. Set up environment variables
  `

  Open http://localhost:3000 (http://localhost:3000) in your browser to see the application.

  ## 📄 API Documentation

  The API is documented using the OpenAPI specification. The openapi.yaml file in the backend directory contains the full API 
  definition. You can use tools like Swagger Editor (https://editor.swagger.io/) to view and interact with the API documentation.

  ## 🤝 Contributing

  Contributions are welcome! Please feel free to submit a pull request or open an issue for any bugs or feature requests.

  ## 📜 License

  This project is licensed under the MIT License. See the LICENSE file for details.
