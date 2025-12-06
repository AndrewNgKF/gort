package generator

func createViewFiles(name string) error {
	layout := `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>{{ .Title }}</title>
    <script src="https://cdn.tailwindcss.com"></script>
    <style>
        body {
            font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, "Helvetica Neue", Arial, sans-serif;
        }
    </style>
</head>
<body class="bg-white min-h-screen">
    <main class="container mx-auto px-6 py-8">
        {{ template "content" . }}
    </main>
</body>
</html>
`

	homeIndex := `{{ define "content" }}
<div class="max-w-4xl mx-auto">
    <!-- Hero Section with Globe Illustration -->
    <div class="text-center mb-16 pt-8">
        <!-- SVG Globe with people holding hands around it -->
        <div class="mb-8">
            <svg viewBox="0 0 400 300" class="w-full max-w-md mx-auto">
                <!-- Globe -->
                <circle cx="200" cy="150" r="80" fill="#E8F4FF" stroke="#6B7280" stroke-width="2"/>
                
                <!-- Latitude lines -->
                <ellipse cx="200" cy="150" rx="80" ry="20" fill="none" stroke="#D1D5DB" stroke-width="1"/>
                <ellipse cx="200" cy="150" rx="80" ry="40" fill="none" stroke="#D1D5DB" stroke-width="1"/>
                
                <!-- Longitude lines -->
                <ellipse cx="200" cy="150" rx="20" ry="80" fill="none" stroke="#D1D5DB" stroke-width="1"/>
                <ellipse cx="200" cy="150" rx="40" ry="80" fill="none" stroke="#D1D5DB" stroke-width="1"/>
                <line x1="200" y1="70" x2="200" y2="230" stroke="#D1D5DB" stroke-width="1"/>
                
                <!-- Continents (simplified) -->
                <path d="M 160 120 Q 170 110 180 115 Q 190 120 185 130 Q 180 135 170 132 Q 165 130 160 120" fill="#9CA3AF"/>
                <path d="M 210 140 Q 225 135 235 145 Q 240 155 230 165 Q 220 168 215 160 Q 210 150 210 140" fill="#9CA3AF"/>
                <path d="M 180 165 Q 190 160 200 170 Q 195 180 185 178 Q 178 172 180 165" fill="#9CA3AF"/>
                
                <!-- People holding hands around the globe (8 people) -->
                <!-- Person 1 (top) -->
                <circle cx="200" cy="50" r="8" fill="#8B5CF6"/>
                <line x1="200" y1="58" x2="200" y2="70" stroke="#8B5CF6" stroke-width="3"/>
                <line x1="200" y1="62" x2="185" y2="55" stroke="#8B5CF6" stroke-width="3"/>
                <line x1="200" y1="62" x2="215" y2="55" stroke="#8B5CF6" stroke-width="3"/>
                
                <!-- Person 2 (top-right) -->
                <circle cx="257" cy="93" r="8" fill="#EC4899"/>
                <line x1="257" y1="101" x2="257" y2="113" stroke="#EC4899" stroke-width="3"/>
                <line x1="257" y1="105" x2="242" y2="98" stroke="#EC4899" stroke-width="3"/>
                <line x1="257" y1="105" x2="270" y2="110" stroke="#EC4899" stroke-width="3"/>
                
                <!-- Person 3 (right) -->
                <circle cx="290" cy="150" r="8" fill="#6366F1"/>
                <line x1="290" y1="158" x2="290" y2="170" stroke="#6366F1" stroke-width="3"/>
                <line x1="290" y1="162" x2="280" y2="155" stroke="#6366F1" stroke-width="3"/>
                <line x1="290" y1="162" x2="303" y2="162" stroke="#6366F1" stroke-width="3"/>
                
                <!-- Person 4 (bottom-right) -->
                <circle cx="257" cy="207" r="8" fill="#10B981"/>
                <line x1="257" y1="215" x2="257" y2="227" stroke="#10B981" stroke-width="3"/>
                <line x1="257" y1="219" x2="242" y2="226" stroke="#10B981" stroke-width="3"/>
                <line x1="257" y1="219" x2="270" y2="215" stroke="#10B981" stroke-width="3"/>
                
                <!-- Person 5 (bottom) -->
                <circle cx="200" cy="250" r="8" fill="#F59E0B"/>
                <line x1="200" y1="258" x2="200" y2="270" stroke="#F59E0B" stroke-width="3"/>
                <line x1="200" y1="262" x2="185" y2="269" stroke="#F59E0B" stroke-width="3"/>
                <line x1="200" y1="262" x2="215" y2="269" stroke="#F59E0B" stroke-width="3"/>
                
                <!-- Person 6 (bottom-left) -->
                <circle cx="143" cy="207" r="8" fill="#EF4444"/>
                <line x1="143" y1="215" x2="143" y2="227" stroke="#EF4444" stroke-width="3"/>
                <line x1="143" y1="219" x2="130" y2="215" stroke="#EF4444" stroke-width="3"/>
                <line x1="143" y1="219" x2="158" y2="226" stroke="#EF4444" stroke-width="3"/>
                
                <!-- Person 7 (left) -->
                <circle cx="110" cy="150" r="8" fill="#06B6D4"/>
                <line x1="110" y1="158" x2="110" y2="170" stroke="#06B6D4" stroke-width="3"/>
                <line x1="110" y1="162" x2="97" y2="162" stroke="#06B6D4" stroke-width="3"/>
                <line x1="110" y1="162" x2="120" y2="155" stroke="#06B6D4" stroke-width="3"/>
                
                <!-- Person 8 (top-left) -->
                <circle cx="143" cy="93" r="8" fill="#A855F7"/>
                <line x1="143" y1="101" x2="143" y2="113" stroke="#A855F7" stroke-width="3"/>
                <line x1="143" y1="105" x2="130" y2="110" stroke="#A855F7" stroke-width="3"/>
                <line x1="143" y1="105" x2="158" y2="98" stroke="#A855F7" stroke-width="3"/>
                
                <!-- Hands connecting (circle around) -->
                <path d="M 200 70 Q 240 75 257 93" fill="none" stroke="#D1D5DB" stroke-width="2" stroke-dasharray="3,3"/>
                <path d="M 257 113 Q 275 130 280 150" fill="none" stroke="#D1D5DB" stroke-width="2" stroke-dasharray="3,3"/>
                <path d="M 290 170 Q 280 190 257 207" fill="none" stroke="#D1D5DB" stroke-width="2" stroke-dasharray="3,3"/>
                <path d="M 257 227 Q 230 240 200 250" fill="none" stroke="#D1D5DB" stroke-width="2" stroke-dasharray="3,3"/>
                <path d="M 185 250 Q 165 240 143 227" fill="none" stroke="#D1D5DB" stroke-width="2" stroke-dasharray="3,3"/>
                <path d="M 143 207 Q 120 190 110 170" fill="none" stroke="#D1D5DB" stroke-width="2" stroke-dasharray="3,3"/>
                <path d="M 110 150 Q 120 120 143 113" fill="none" stroke="#D1D5DB" stroke-width="2" stroke-dasharray="3,3"/>
                <path d="M 143 93 Q 165 75 185 70" fill="none" stroke="#D1D5DB" stroke-width="2" stroke-dasharray="3,3"/>
            </svg>
        </div>
        
        <h1 class="text-5xl font-light text-gray-800 mb-4">
            {{ .Title }}
        </h1>
        <p class="text-xl text-gray-600 font-light mb-2">{{ .Message }}</p>
        <p class="text-base text-gray-500">A Rails-inspired web framework for Go</p>
    </div>

    <!-- Simple Info Sections -->
    <div class="space-y-12 mb-16">
        <div class="text-center">
            <h2 class="text-2xl font-light text-gray-800 mb-4">About your application's environment</h2>
            <div class="bg-white rounded-lg shadow-sm border border-gray-200 p-6 max-w-2xl mx-auto">
                <div class="grid grid-cols-2 gap-4 text-left">
                    <div class="border-b border-gray-100 pb-2">
                        <span class="text-gray-600 text-sm">Go version</span>
                    </div>
                    <div class="border-b border-gray-100 pb-2">
                        <span class="text-gray-800 text-sm font-mono">{{ .GoVersion }}</span>
                    </div>
                    <div class="border-b border-gray-100 pb-2">
                        <span class="text-gray-600 text-sm">Gort version</span>
                    </div>
                    <div class="border-b border-gray-100 pb-2">
                        <span class="text-gray-800 text-sm font-mono">{{ .GortVersion }}</span>
                    </div>
                    <div class="pb-2">
                        <span class="text-gray-600 text-sm">Environment</span>
                    </div>
                    <div class="pb-2">
                        <span class="text-gray-800 text-sm font-mono">development</span>
                    </div>
                </div>
            </div>
        </div>

        <div class="text-center">
            <h2 class="text-2xl font-light text-gray-800 mb-4">Getting Started</h2>
            <div class="bg-white rounded-lg shadow-sm border border-gray-200 p-6 max-w-2xl mx-auto text-left">
                <p class="text-gray-600 mb-4">Here are some common commands to help you build your application:</p>
                <div class="space-y-3">
                    <div class="bg-gray-50 rounded px-4 py-2 font-mono text-sm text-gray-700">
                        <span class="text-gray-500">$</span> gort db:setup
                    </div>
                    <p class="text-xs text-gray-500 px-4 -mt-2 mb-2">Create database and run migrations</p>
                    <div class="bg-gray-50 rounded px-4 py-2 font-mono text-sm text-gray-700">
                        <span class="text-gray-500">$</span> gort generate scaffold Post title:string body:text
                    </div>
                    <p class="text-xs text-gray-500 px-4 -mt-2 mb-2">Generate a complete CRUD resource</p>
                    <div class="bg-gray-50 rounded px-4 py-2 font-mono text-sm text-gray-700">
                        <span class="text-gray-500">$</span> gort server
                    </div>
                    <p class="text-xs text-gray-500 px-4 -mt-2">Start your application server</p>
                </div>
            </div>
        </div>
    </div>

    <!-- Footer Links -->
    <div class="text-center pb-12">
        <div class="flex gap-6 justify-center text-sm">
            <a href="https://github.com/AndrewNgKF/gort" target="_blank" class="text-gray-600 hover:text-gray-900 transition">
                View on GitHub
            </a>
            <span class="text-gray-300">•</span>
            <a href="https://pkg.go.dev/github.com/AndrewNgKF/gort" target="_blank" class="text-gray-600 hover:text-gray-900 transition">
                Documentation
            </a>
            <span class="text-gray-300">•</span>
            <a href="https://github.com/AndrewNgKF/gort/blob/main/README.md" target="_blank" class="text-gray-600 hover:text-gray-900 transition">
                Guides
            </a>
        </div>
    </div>
</div>
{{ end }}
`

	if err := writeFile(name, "app/views/layouts/application.html", layout); err != nil {
		return err
	}
	return writeFile(name, "app/views/home/index.html", homeIndex)
}

func createSeedsFile(name string) error {
	content := `package main

import (
	"fmt"
)

func main() {
	fmt.Println("Seeding database...")
	fmt.Println("✅ Seeding complete!")
}
`
	return writeFile(name, "db/seeds.go", content)
}

func createCSSFile(name string) error {
	content := `/* Custom styles for your Gort application */
/* Tailwind CSS is loaded via CDN in the layout */

/* Add your custom CSS here */
`
	return writeFile(name, "public/assets/css/application.css", content)
}
