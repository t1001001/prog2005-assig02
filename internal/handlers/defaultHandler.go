package handlers

import (
	"net/http"

	"cloud.google.com/go/firestore"
)

func DefaultHandler(w http.ResponseWriter, r *http.Request, client *firestore.Client) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)

	html := `
<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8" />
  <meta name="viewport" content="width=device-width, initial-scale=1.0" />
  <title>Countries Dashboards Service</title>
  <link href="https://fonts.googleapis.com/css2?family=Space+Grotesk:wght@400;700&display=swap" rel="stylesheet">
  <style>
    body, html {
      margin: 0;
      padding: 0;
      height: 100%;
      font-family: 'Space Grotesk', sans-serif;
      color: white;
      overflow: hidden;
      position: relative;
      z-index: 1;
    }

    body::before {
      content: "";
      position: fixed;
      top: 0;
      left: 0;
      width: 100%;
      height: 100%;
      background: url('/images/earth-view.jpg') no-repeat center center fixed;
      background-size: cover;
      transform: rotate(180deg);
      z-index: -1;
      opacity: 1;
    }

    nav {
      position: absolute;
      top: 0;
      right: 0;
      padding: 30px 50px;
      z-index: 2;
    }

    nav ul {
      list-style: none;
      display: flex;
      gap: 30px;
      margin: 0;
      padding: 0;
    }

    nav ul li a {
      color: white;
      text-decoration: none;
      font-weight: 600;
      letter-spacing: 1px;
      transition: color 0.3s ease;
      user-select: none;
      -webkit-user-select: none;
      -ms-user-select: none;
      cursor: default;
    }

    nav ul li a:hover {
      color: #00bfff;
    }

    .content {
      height: 100vh;
      display: flex;
      align-items: center;
      justify-content: flex-start;
      padding-left: 8%;
      max-width: 700px;
      animation: fadeInUp 1.2s ease forwards;
      opacity: 0;
    }

    @keyframes fadeInUp {
      0% {
        opacity: 0;
        transform: translateY(30px);
      }
      100% {
        opacity: 1;
        transform: translateY(0);
      }
    }

    .content h2 {
      font-size: 1em;
      text-transform: uppercase;
      letter-spacing: 3px;
      color: #ccc;
      margin-bottom: 10px;
    }

    .content h1 {
      font-size: 4em;
      margin-bottom: 20px;
      font-weight: 700;
      letter-spacing: 2px;
    }

    .content p {
      font-size: 1.1em;
      line-height: 1.6;
      color: #ddd;
    }

    @media (max-width: 768px) {
      .content {
        padding: 0 5%;
        justify-content: center;
        text-align: center;
      }
      .content h1 {
        font-size: 2.5em;
      }
    }
  </style>
</head>
<body>
  <nav>
    <ul>
      <li><a href="/dashboard/v1/registrations/">Registrations</a></li>
      <li><a href="/dashboard/v1/status/">Status</a></li>
      <li><a href="/dashboard/v1/dashboards/">Dashboards</a></li>
      <li><a href="/dashboard/v1/notifications/">Notifications</a></li>
    </ul>
  </nav>

  <div class="content">
    <div>
      <h2>So, you want to dashboard the world?</h2>
      <h1>Countries Dashboards Service</h1>
      <p>
        Let’s face it: if you want to manage data, you might as well truly harness the power of dynamic dashboards — not hover around basic visualizations.
        Well, sit back and explore. We’ll let you deploy dashboards that are connected, reactive, and ready for real-world data — with persistence,
        notifications, and infrastructure that scale.
      </p>
    </div>
  </div>
</body>
</html>
`
	w.Write([]byte(html))
}
