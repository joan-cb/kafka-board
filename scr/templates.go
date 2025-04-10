package main

var homeTemplate string = `<!DOCTYPE html>
<html>
<head>
    <title>Schema Registry Dashboard</title>
    <style>
        :root {
            --primary-color: #4a90e2;
            --primary-dark: #357abd;
            --primary-light: #e8f2f9;
            --secondary-color: #9b59b6;
            --secondary-light: #f5eef8;
            --text-primary: #2c3e50;
            --text-secondary: #546e7a;
            --background-start: #f0f4f8;
            --background-end: #e8f2f9;
            --card-background: #ffffff;
            --shadow-color: rgba(0, 0, 0, 0.1);
            --transition-speed: 0.3s;
        }

        body {
            font-family: 'Segoe UI', Arial, sans-serif;
            max-width: 1200px;
            margin: 0 auto;
            padding: 20px;
            background: #cce6ff;
            color: var(--text-primary);
            line-height: 1.6;
            position: relative;
            min-height: 100vh;
            padding-bottom: 80px;
        }

        h1 {
            color: var(--primary-color);
            text-align: center;
            margin-bottom: 30px;
            font-size: 2.5em;
            text-shadow: 2px 2px 4px var(--shadow-color);
            position: relative;
            padding-bottom: 15px;
            background: linear-gradient(135deg, var(--primary-color), var(--secondary-color));
            -webkit-background-clip: text;
            -webkit-text-fill-color: transparent;
            font-weight: 800;
            letter-spacing: -0.5px;
        }

        h1::after {
            content: '';
            position: absolute;
            bottom: 0;
            left: 50%;
            transform: translateX(-50%);
            width: 150px;
            height: 4px;
            background: linear-gradient(to right, var(--primary-color), var(--secondary-color));
            border-radius: 2px;
            box-shadow: 0 2px 4px var(--shadow-color);
        }

        .header-container {
            text-align: center;
            margin-bottom: 30px;
            padding: 20px;
            background: rgba(255, 255, 255, 0.9);
            border-radius: 15px;
            box-shadow: 0 4px 6px var(--shadow-color);
            backdrop-filter: blur(10px);
        }

        .header-image {
            max-width: 100%;
            height: auto;
            display: block;
            margin: 0 auto;
            width: 800px;
            margin-bottom: 10px;
        }

        .header-subtitle {
            color: var(--text-secondary);
            font-size: 1.1em;
            margin-top: 10px;
            font-weight: 500;
        }

        .header-stats {
            display: flex;
            justify-content: center;
            gap: 20px;
            margin-top: 20px;
            flex-wrap: wrap;
        }

        .stat-item {
            background: var(--primary-light);
            padding: 10px 20px;
            border-radius: 20px;
            color: var(--primary-dark);
            font-weight: 600;
            display: flex;
            align-items: center;
            gap: 8px;
            box-shadow: 0 2px 4px var(--shadow-color);
            transition: all var(--transition-speed) ease;
        }

        .slack-button {
            background: none;
            border: none;
            padding: 0;
            cursor: pointer;
            transition: transform var(--transition-speed) ease;
        }

        .slack-button:hover {
            transform: scale(1.05);
        }

        .slack-button img {
            height: 40px;
            width: auto;
        }

        .github-button {
            background: none;
            border: none;
            padding: 0;
            cursor: pointer;
            transition: transform var(--transition-speed) ease;
        }

        .github-button:hover {
            transform: scale(1.05);
        }

        .github-button img {
            height: 40px;
            width: auto;
        }

        .stat-item:hover {
            transform: translateY(-2px);
            box-shadow: 0 4px 8px var(--shadow-color);
        }

        .stat-icon {
            font-size: 1.2em;
            width: 24px;
            height: 24px;
            object-fit: contain;
        }

        .subject-card {
            background-color: var(--card-background);
            border: none;
            padding: 25px;
            margin: 20px 0;
            border-radius: 12px;
            box-shadow: 0 4px 6px var(--shadow-color);
            transition: all var(--transition-speed) ease;
            position: relative;
            overflow: hidden;
            backdrop-filter: blur(10px);
            background: rgba(255, 255, 255, 0.95);
        }

        .subject-card::before {
            content: '';
            position: absolute;
            top: 0;
            left: 0;
            width: 4px;
            height: 100%;
            background: linear-gradient(to bottom, var(--primary-color), var(--secondary-color));
            opacity: 0;
            transition: opacity var(--transition-speed) ease;
        }

        .subject-card:hover {
            transform: translateY(-5px);
            box-shadow: 0 8px 15px var(--shadow-color);
        }

        .subject-card:hover::before {
            opacity: 1;
        }

        .subject-name {
            font-weight: 600;
            color: var(--primary-color);
            font-size: 1.4em;
            margin-bottom: 15px;
            display: flex;
            align-items: center;
        }

        .subject-emoji {
            margin-right: 10px;
            font-size: 1.2em;
        }

        .property {
            margin: 12px 0;
            color: var(--text-secondary);
            display: flex;
            align-items: center;
            position: relative;
        }

        .property-label {
            font-weight: 600;
            color: var(--primary-dark);
            min-width: 150px;
            padding: 5px 10px;
            background-color: var(--primary-light);
            border-radius: 4px;
            margin-right: 10px;
        }

        .property-value {
            display: flex;
            align-items: center;
            gap: 8px;
        }

        .icon-badge {
            display: inline-flex;
            align-items: center;
            padding: 4px 8px;
            border-radius: 12px;
            font-size: 0.9em;
            font-weight: 600;
            gap: 4px;
            box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
        }

        .icon-badge-true {
            background-color: #e8f5e9;
            color: #2ecc71;
            border: 1px solid #2ecc71;
        }

        .icon-badge-false {
            background-color: #ffebee;
            color: #e74c3c;
            border: 1px solid #e74c3c;
        }

        .icon-badge-backward {
            background-color: #fff3e0;
            color: #f39c12;
            border: 1px solid #f39c12;
        }

        .icon-badge-forward {
            background-color: #e3f2fd;
            color: #3498db;
            border: 1px solid #3498db;
        }

        .icon-badge-full {
            background-color: #f3e5f5;
            color: #9b59b6;
            border: 1px solid #9b59b6;
        }

        .icon-badge-type {
            background-color: #f3e5f5;
            color: #9b59b6;
            border: 1px solid #9b59b6;
        }

        .icon-badge-none {
            background-color: #eceff1;
            color: #95a5a6;
            border: 1px solid #bdc3c7;
        }
        
        .icon-badge-warning {
            background-color: #fff3e0;
            color: #f39c12;
            border: 1px solid #f39c12;
        }
        
        .icon-badge-id {
            background-color: #e3f2fd;
            color: #3498db;
            border: 1px solid #3498db;
        }
        
        .icon-badge-version {
            background-color: #e3f2fd;
            color: #3498db;
            border: 1px solid #3498db;
        }
        
        .icon-badge-subject {
            background-color: #f3e5f5;
            color: #9b59b6;
            display: inline-flex;
            align-items: center;
            padding: 4px 8px;
            border-radius: 12px;
            font-size: 0.9em;
            font-weight: 600;
            gap: 4px;
            box-shadow: 0 2px 4px rgba(155, 89, 182, 0.2);
            border: 1px solid #9b59b6;
        }

        .icon {
            font-size: 1.2em;
        }

        .icon-true {
            color: #2ecc71;
        }

        .icon-false {
            color: #e74c3c;
        }

        .icon-backward {
            color: #f39c12;
        }

        .icon-forward {
            color: #3498db;
        }

        .icon-full {
            color: #9b59b6;
        }

        .icon-none {
            color: #95a5a6;
        }

        .alias-list {
            display: flex;
            flex-wrap: wrap;
            gap: 8px;
        }

        .alias-tag {
            background-color: var(--secondary-light);
            padding: 4px 8px;
            border-radius: 12px;
            color: var(--secondary-color);
            font-size: 0.9em;
            font-weight: 600;
            display: inline-flex;
            align-items: center;
            gap: 4px;
            transition: all var(--transition-speed) ease;
        }

        .alias-tag:hover {
            transform: scale(1.05);
            box-shadow: 0 2px 4px var(--shadow-color);
        }

        button {
            background: linear-gradient(135deg, var(--primary-color), var(--primary-dark));
            color: white;
            border: none;
            padding: 12px 24px;
            border-radius: 25px;
            cursor: pointer;
            font-size: 1em;
            font-weight: 600;
            transition: all var(--transition-speed) ease;
            box-shadow: 0 2px 4px var(--shadow-color);
            margin: 10px 0;
            display: block;
            margin-left: auto;
        }

        button:hover {
            transform: translateY(-2px);
            box-shadow: 0 4px 8px var(--shadow-color);
        }

        button:active {
            transform: translateY(0);
        }

        .search-container {
            margin-bottom: 30px;
            text-align: center;
            position: relative;
            display: flex;
            flex-direction: column;
            align-items: center;
            width: 100%;
        }

        .subject-counter {
            background: var(--primary-light);
            padding: 10px 20px;
            border-radius: 20px;
            color: var(--primary-dark);
            font-weight: 600;
            display: inline-flex;
            align-items: center;
            gap: 8px;
            box-shadow: 0 2px 4px var(--shadow-color);
            margin-bottom: 15px;
        }

        .search-input {
            width: 80%;
            padding: 15px 20px;
            font-size: 16px;
            border: 2px solid var(--primary-light);
            border-radius: 25px;
            outline: none;
            transition: all var(--transition-speed) ease;
            background-color: var(--card-background);
            box-shadow: 0 2px 4px var(--shadow-color);
        }

        .search-input:focus {
            border-color: var(--primary-color);
            box-shadow: 0 4px 12px rgba(33, 150, 243, 0.2);
            width: 85%;
        }

        .search-input::placeholder {
            color: var(--text-secondary);
        }

        .no-results {
            text-align: center;
            padding: 40px 20px;
            color: var(--text-secondary);
            font-family: monospace;
            white-space: pre;
            line-height: 1.2;
            background: var(--card-background);
            border-radius: 12px;
            box-shadow: 0 4px 6px var(--shadow-color);
            margin: 20px auto;
            max-width: 600px;
        }

        .no-results-text {
            margin-top: 20px;
            font-family: 'Segoe UI', Arial, sans-serif;
            color: var(--primary-color);
            font-size: 1.2em;
        }

        .hidden {
            display: none;
            opacity: 0;
            transform: translateY(20px);
            transition: all var(--transition-speed) ease;
        }

        @media (max-width: 768px) {
            body {
                padding: 15px;
            }

            .subject-card {
                padding: 20px;
            }

            .property {
                flex-direction: column;
                align-items: flex-start;
            }

            .property-label {
                margin-bottom: 5px;
            }

            .search-input {
                width: 90%;
            }

            .search-input:focus {
                width: 95%;
            }
        }

        .footer {
            position: fixed;
            bottom: 0;
            left: 0;
            right: 0;
            text-align: center;
            padding: 10px;
            background: rgba(255, 255, 255, 0.9);
            backdrop-filter: blur(10px);
            z-index: 100;
        }

        .footer-image {
            max-width: 100%;
            height: auto;
            display: block;
            margin: 0 auto;
            width: 300px;
        }

        .back-button {
            display: inline-block;
            background: linear-gradient(135deg, var(--primary-color), var(--primary-dark));
            color: white;
            border: none;
            padding: 12px 24px;
            border-radius: 25px;
            cursor: pointer;
            font-size: 1em;
            font-weight: 600;
            transition: all var(--transition-speed) ease;
            margin-bottom: 15px;
            box-shadow: 0 2px 4px var(--shadow-color);
            text-decoration: none;
        }

        .test-button {
            display: inline-block;
            background: linear-gradient(135deg, var(--primary-color), var(--primary-dark));
            color: white;
            border: none;
            padding: 12px 24px;
            border-radius: 25px;
            cursor: pointer;
            font-size: 1em;
            font-weight: 600;
            transition: all var(--transition-speed) ease;
            margin-bottom: 15px;
            margin-left: 15px;
            box-shadow: 0 2px 4px var(--shadow-color);
            text-decoration: none;
        }

        .icon-badge-text {
            margin-left: 4px;
        }
    </style>
</head>
<body>
    <div class="header-container">
        <img src="/static/header.png" class="header-image" alt="Header">
        <div class="header-stats">
            <a href="https://slack.com" target="_blank" class="slack-button">
                <img src="/static/slack-channel.png" class="slack-image" alt="Slack">
            </a>
            <a href="https://www.lemonde.fr" target="_blank" class="github-button">
                <img src="/static/github.png" class="github-image" alt="GitHub">
            </a>
        </div>
    </div>
    <div class="search-container">
        <div class="subject-counter">
            <span class="stat-icon">📊</span>
            <span>{{len .Configs}} Subjects</span>
        </div>
        <input type="text" id="searchInput" class="search-input" placeholder="Search subjects..." onkeyup="filterSubjects()">
    </div>

    <!-- Global Config Card -->
    <div id="globalConfig" class="subject-card global-config">
        <div class="subject-name">
            <span class="subject-emoji">🌍</span>
            {{.GlobalConfig.Name}}
        </div>
        <div class="property">
            <span class="property-label">Normalize:</span>
            <div class="property-value">
                {{with .GlobalConfig.Normalize}}
                    {{if .}}
                        <span class="icon-badge icon-badge-true">True</span>
                    {{else}}
                        <span class="icon-badge icon-badge-false">False</span>
                    {{end}}
                {{else}}
                    <span class="icon-badge icon-badge-none">Not Set</span>
                {{end}}
            </div>
        </div>
        <div class="property">
            <span class="property-label">Compatibility Level:</span>
            <div class="property-value">
                {{if eq .GlobalConfig.CompatibilityLevel "BACKWARD"}}
                    <span class="icon-badge icon-badge-backward">Backward</span>
                {{else if eq .GlobalConfig.CompatibilityLevel "FORWARD"}}
                    <span class="icon-badge icon-badge-forward">Forward</span>
                {{else if eq .GlobalConfig.CompatibilityLevel "FULL"}}
                    <span class="icon-badge icon-badge-full">Full</span>
                {{else}}
                    <span class="icon-badge icon-badge-none">{{.GlobalConfig.CompatibilityLevel}}</span>
                {{end}}
            </div>
        </div>
        <div class="property">
            <span class="property-label">Aliases:</span>
            <div class="alias-list">
                <span class="alias-tag">{{.GlobalConfig.Alias}}</span>
            </div>
        </div>
    </div>

    <!-- Subject Configs -->
    <div id="subjectConfigs">
        {{range .Configs}}
        <div class="subject-card" data-name="{{.GetName}}">
            <div class="subject-name">
                <span class="subject-emoji"></span>
                {{.GetName}}
            </div>
            {{if eq (printf "%T" .) "main.SubjectConfig"}}
            <div class="property">
                <span class="property-label">Normalize:</span>
                <div class="property-value">
                    {{with .Normalize}}
                        {{if .}}
                            <span class="icon-badge icon-badge-true">True</span>
                        {{else}}
                            <span class="icon-badge icon-badge-false">False</span>
                        {{end}}
                    {{else}}
                        <span class="icon-badge icon-badge-none">Not Set</span>
                    {{end}}
                </div>
            </div>
            <div class="property">
                <span class="property-label">Compatibility Level:</span>
                <div class="property-value">
                    {{if eq .CompatibilityLevel "BACKWARD"}}
                        <span class="icon-badge icon-badge-backward">Backward</span>
                    {{else if eq .CompatibilityLevel "FORWARD"}}
                        <span class="icon-badge icon-badge-forward">Forward</span>
                    {{else if eq .CompatibilityLevel "FULL"}}
                        <span class="icon-badge icon-badge-full">Full</span>
                    {{else}}
                        <span class="icon-badge icon-badge-none">{{.CompatibilityLevel}}</span>
                    {{end}}
                </div>
            </div>
            <div class="property">
                <span class="property-label">Aliases:</span>
                <div class="alias-list">
                    <span class="alias-tag">{{.Alias}}</span>
                </div>
            </div>
            {{else}}
            <div class="property">
                <span class="property-label">Status:</span>
                <div class="property-value">
                    <span class="icon-badge icon-badge-none">Using Global Default</span>
                </div>
            </div>
            {{end}}
            <button onclick="viewSchema('{{.GetName}}')">View Schema</button>
        </div>
        {{end}}
    </div>

    <div id="no-results" class="no-results hidden">
        <pre>
    ╔══════════════════════════════════════╗
    ║  ¯\_(ツ)_/¯                          ║
    ║                                     ║
    ║  No subjects found!                 ║
    ║                                     ║
    ║  Try a different search term        ║
    ║                                     ║
    ║  (｡•́︿•̀｡)                           ║
    ╚══════════════════════════════════════╝
        </pre>
    </div>

    <div class="footer">
        <img src="/static/footer.png" class="footer-image" alt="Footer">
    </div>

    <script>
        const subjectEmojis = [
            '📋', '📝', '📄', '📑', '📒', '📓', '📔', '📕', '📗', '📘',
            '📙', '📚', '📖', '📗', '📘', '📙', '📚', '📖', '📗', '📘',
            '📙', '📚', '📖', '📗', '📘', '📙', '📚', '📖', '📗', '📘',
            '📙', '📚', '📖', '📗', '📘', '📙', '📚', '📖', '📗', '📘',
            '📙', '📚', '📖', '📗', '📘', '📙', '📚', '📖', '📗', '📘',
            '📙', '📚', '📖', '📗', '📘', '📙', '📚', '📖', '📗', '📘',
            '📙', '📚', '📖', '📗', '📘', '📙', '📚', '📖', '📗', '📘',
            '📙', '📚', '📖', '📗', '📘', '📙', '📚', '📖', '📗', '📘',
            '📙', '📚', '📖', '📗', '📘', '📙', '📚', '📖', '📗', '📘',
            '📙', '📚', '📖', '📗', '📘', '📙', '📚', '📖', '📗', '📘'
        ];

        function getRandomEmoji() {
            const randomIndex = Math.floor(Math.random() * subjectEmojis.length);
            return subjectEmojis[randomIndex];
        }

        function viewSchema(topicName) {
            window.location.href = '/schema/?topic=' + encodeURIComponent(topicName);
        }

        function filterSubjects() {
            const input = document.getElementById('searchInput');
            const filter = input.value.toUpperCase();
            const subjectConfigs = document.getElementById('subjectConfigs');
            const globalConfig = document.getElementById('globalConfig');
            const cards = subjectConfigs.getElementsByClassName('subject-card');
            const noResults = document.getElementById('no-results');
            let visibleCount = 0;
            
            // Hide subjects container by default
            subjectConfigs.style.display = 'none';
            noResults.classList.add('hidden');
            
            if (filter === '') {
                // If no search term, only show global config
                globalConfig.style.display = 'block';
                return;
            }

            // If there's a search term, hide global config and show subjects container
            globalConfig.style.display = 'none';
            subjectConfigs.style.display = 'block';
            
            for (let i = 0; i < cards.length; i++) {
                const subjectName = cards[i].querySelector('.subject-name');
                if (subjectName) {
                    const txtValue = subjectName.textContent || subjectName.innerText;
                    if (txtValue.toUpperCase().indexOf(filter) > -1) {
                        cards[i].style.display = '';
                        visibleCount++;
                    } else {
                        cards[i].style.display = 'none';
                    }
                }
            }

            if (visibleCount === 0) {
                noResults.classList.remove('hidden');
            } else {
                noResults.classList.add('hidden');
            }
        }

        // Call filterSubjects on page load to set initial state
        document.addEventListener('DOMContentLoaded', function() {
            filterSubjects();
            const subjectCards = document.querySelectorAll('.subject-card:not(.global-config)');
            subjectCards.forEach(card => {
                const emojiSpan = card.querySelector('.subject-emoji');
                if (emojiSpan) {
                    emojiSpan.textContent = getRandomEmoji();
                }
            });
        });
    </script>
</body>
</html>`

var schemaTemplate string = `<!DOCTYPE html>
<html>
<head>
    <title>Schema Viewer - {{.SubjectName}}</title>
    <style>
        :root {
            --primary-color: #4a90e2;
            --primary-dark: #357abd;
            --primary-light: #e8f2f9;
            --secondary-color: #9b59b6;
            --secondary-light: #f5eef8;
            --text-primary: #2c3e50;
            --text-secondary: #546e7a;
            --background-start: #f0f4f8;
            --background-end: #e8f2f9;
            --card-background: #ffffff;
            --shadow-color: rgba(0, 0, 0, 0.1);
            --transition-speed: 0.3s;
        }

        body {
            font-family: 'Segoe UI', Arial, sans-serif;
            max-width: 1200px;
            margin: 0 auto;
            padding: 20px;
            background: #cce6ff;
            color: var(--text-primary);
            line-height: 1.6;
            min-height: 100vh;
            position: relative;
            padding-bottom: 80px;
        }

        .header-container {
            text-align: center;
            margin-bottom: 30px;
            padding: 20px;
            background: rgba(255, 255, 255, 0.9);
            border-radius: 15px;
            box-shadow: 0 4px 6px var(--shadow-color);
            backdrop-filter: blur(10px);
            position: relative;
        }

        .header-image {
            max-width: 100%;
            height: auto;
            display: block;
            margin: 0 auto;
            width: 800px;
            margin-bottom: 10px;
        }

        .header-stats {
            display: flex;
            justify-content: center;
            gap: 20px;
            margin-top: 20px;
            flex-wrap: wrap;
        }

        .slack-button {
            background: none;
            border: none;
            padding: 0;
            cursor: pointer;
            transition: transform var(--transition-speed) ease;
        }

        .slack-button:hover {
            transform: scale(1.05);
        }

        .slack-button img {
            height: 40px;
            width: auto;
        }

        .github-button {
            background: none;
            border: none;
            padding: 0;
            cursor: pointer;
            transition: transform var(--transition-speed) ease;
        }

        .github-button:hover {
            transform: scale(1.05);
        }

        .github-button img {
            height: 40px;
            width: auto;
        }

        h1 {
            color: var(--primary-color);
            text-align: center;
            margin: 20px 0;
            font-size: 2.5em;
            text-shadow: 2px 2px 4px var(--shadow-color);
            position: relative;
            padding-bottom: 15px;
            background: linear-gradient(135deg, var(--primary-color), var(--secondary-color));
            -webkit-background-clip: text;
            -webkit-text-fill-color: transparent;
            font-weight: 800;
            letter-spacing: -0.5px;
        }

        h1::after {
            content: '';
            position: absolute;
            bottom: 0;
            left: 50%;
            transform: translateX(-50%);
            width: 150px;
            height: 4px;
            background: linear-gradient(to right, var(--primary-color), var(--secondary-color));
            border-radius: 2px;
            box-shadow: 0 2px 4px var(--shadow-color);
        }

        .button-container {
            display: flex;
            justify-content: space-between;
            margin-bottom: 20px;
        }

        .back-button {
            background: linear-gradient(135deg, var(--primary-color), var(--primary-dark));
            color: white;
            border: none;
            padding: 12px 24px;
            border-radius: 25px;
            cursor: pointer;
            font-size: 1em;
            font-weight: 600;
            transition: all var(--transition-speed) ease;
            box-shadow: 0 2px 4px var(--shadow-color);
            text-decoration: none;
            position: absolute;
            bottom: 10px;
            left: 50%;
        }

        .test-button {
            background: linear-gradient(135deg, var(--primary-color), var(--primary-dark));
            color: white;
            border: none;
            padding: 12px 24px;
            border-radius: 25px;
            cursor: pointer;
            font-size: 1em;
            font-weight: 600;
            transition: all var(--transition-speed) ease;
            box-shadow: 0 2px 4px var(--shadow-color);
        }

        .back-button:hover, .test-button:hover {
            transform: translateY(-2px);
            box-shadow: 0 4px 8px var(--shadow-color);
        }

        .back-button:active, .test-button:active {
            transform: translateY(0);
        }

        .schema-card {
            background: white;
            border-radius: 8px;
            padding: 20px;
            margin: 20px 0;
            box-shadow: 0 2px 4px rgba(0,0,0,0.1);
            position: relative;
        }

        .schema-content {
            background: #f8f8f8;
            padding: 15px;
            border-radius: 4px;
            white-space: pre-wrap;
            font-family: monospace;
            overflow-x: auto;
        }

        .schema-content pre {
            margin: 0;
            padding: 0;
            white-space: pre-wrap;
            word-wrap: break-word;
            font-family: 'Consolas', 'Monaco', 'Courier New', monospace;
            font-size: 14px;
            line-height: 1.5;
        }

        .property {
            margin: 12px 0;
            color: var(--text-secondary);
            display: flex;
            align-items: center;
        }

        .property-label {
            font-weight: 600;
            color: var(--primary-dark);
            min-width: 150px;
            padding: 5px 10px;
            background-color: var(--primary-light);
            border-radius: 4px;
            margin-right: 10px;
        }

        .property-value {
            display: flex;
            align-items: center;
            gap: 8px;
        }

        .icon-badge {
            display: inline-flex;
            align-items: center;
            padding: 4px 8px;
            border-radius: 12px;
            font-size: 0.9em;
            font-weight: 600;
            gap: 4px;
            box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
        }

        .icon-badge-version {
            background-color: #e3f2fd;
            color: #3498db;
        }

        .icon-badge-id {
            background-color: #e8f5e9;
            color: #2ecc71;
        }

        .icon-badge-type {
            background-color: #fff3e0;
            color: #f39c12;
        }

        .icon-badge-subject {
            background-color: #f3e5f5;
            color: #9b59b6;
            display: inline-flex;
            align-items: center;
            padding: 4px 8px;
            border-radius: 12px;
            font-size: 0.9em;
            font-weight: 600;
            gap: 4px;
            box-shadow: 0 2px 4px rgba(155, 89, 182, 0.2);
        }

        .footer {
            position: fixed;
            bottom: 0;
            left: 0;
            right: 0;
            text-align: center;
            padding: 10px;
            background: rgba(255, 255, 255, 0.9);
            backdrop-filter: blur(10px);
            z-index: 100;
        }

        .footer-image {
            max-width: 100%;
            height: auto;
            display: block;
            margin: 0 auto;
            width: 300px;
        }
    </style>
</head>
<body>
    <div class="header-container">
        <a href="/" class="back-button">Back to Dashboard</a>
        <img src="/static/header.png" class="header-image" alt="Header">
        <div class="header-stats">
            <a href="https://slack.com" target="_blank" class="slack-button">
                <img src="/static/slack-channel.png" class="slack-image" alt="Slack">
            </a>
            <a href="https://www.lemonde.fr" target="_blank" class="github-button">
                <img src="/static/github.png" class="github-image" alt="GitHub">
            </a>
        </div>
        <h1>{{.SubjectName}}</h1>
    </div>
    {{range .Schemas}}
    <div class="schema-card">
        <button class="test-button" onclick="testSchema('{{$.SubjectName}}', {{.Version}}, {{.Id}})">Test against this schema</button>
        <div class="property">
            <span class="property-label">Version:</span>
            <div class="property-value">
                <span class="icon-badge icon-badge-version">{{.Version}}</span>
            </div>
        </div>
        <div class="property">
            <span class="property-label">ID:</span>
            <div class="property-value">
                <span class="icon-badge icon-badge-id">{{.Id}}</span>
            </div>
        </div>
        <div class="property">
            <span class="property-label">Schema Type:</span>
            <div class="property-value">
                <span class="icon-badge icon-badge-type">{{.SchemaType}}</span>
            </div>
        </div>
        <div class="property">
            <span class="property-label">Schema:</span>
            <div class="schema-content">
                <pre>{{.Schema | formatJSON | html}}</pre>
            </div>
        </div>
    </div>
    {{end}}
    <div class="footer">
        <img src="/static/footer.png" class="footer-image" alt="Footer">
    </div>

    <script>
        function testSchema(subjectName, version, id) {
            window.location.href = '/test-schema/?topic=' + encodeURIComponent(subjectName) + 
                                 '&version=' + encodeURIComponent(version) + 
                                 '&id=' + encodeURIComponent(id);
        }
    </script>
</body>
</html>`

var testSchemaTemplate string = `<!DOCTYPE html>
<html>
<head>
    <title>Test Schema Compatibility</title>
    <style>
        :root {
            --primary-color: #4a90e2;
            --primary-dark: #357abd;
            --primary-light: #e8f2f9;
            --secondary-color: #9b59b6;
            --secondary-light: #f5eef8;
            --text-primary: #2c3e50;
            --text-secondary: #546e7a;
            --background-start: #f0f4f8;
            --background-end: #e8f2f9;
            --card-background: #ffffff;
            --shadow-color: rgba(0, 0, 0, 0.1);
            --transition-speed: 0.3s;
            
            /* Badge Colors */
            --badge-success-bg: #e3f9e5;
            --badge-success-text: #1e8e3e;
            --badge-error-bg: #fdeced;
            --badge-error-text: #d93025;
            --badge-warning-bg: #fff8e1;
            --badge-warning-text: #f57c00;
            --badge-info-bg: #e8f0fe;
            --badge-info-text: #1a73e8;
            --badge-neutral-bg: #f5f5f5;
            --badge-neutral-text: #5f6368;
            --badge-type-bg: #f3e5f5;
            --badge-type-text: #9c27b0;
        }

        body {
            font-family: 'Segoe UI', Arial, sans-serif;
            max-width: 1200px;
            margin: 0 auto;
            padding: 20px;
            background: #cce6ff;
            color: var(--text-primary);
            line-height: 1.6;
            min-height: 100vh;
            position: relative;
            padding-bottom: 80px;
        }

        .header-container {
            text-align: center;
            margin-bottom: 30px;
            padding: 20px;
            background: rgba(255, 255, 255, 0.9);
            border-radius: 15px;
            box-shadow: 0 4px 6px var(--shadow-color);
            backdrop-filter: blur(10px);
            position: relative;
        }

        .header-image {
            max-width: 100%;
            height: auto;
            display: block;
            margin: 0 auto;
            width: 800px;
            margin-bottom: 10px;
        }

        .header-stats {
            display: flex;
            justify-content: center;
            gap: 20px;
            margin-top: 20px;
            flex-wrap: wrap;
        }

        .slack-button {
            background: none;
            border: none;
            padding: 0;
            cursor: pointer;
            transition: transform var(--transition-speed) ease;
        }

        .slack-button:hover {
            transform: scale(1.05);
        }

        .slack-button img {
            height: 40px;
            width: auto;
        }

        .github-button {
            background: none;
            border: none;
            padding: 0;
            cursor: pointer;
            transition: transform var(--transition-speed) ease;
        }

        .github-button:hover {
            transform: scale(1.05);
        }

        .github-button img {
            height: 40px;
            width: auto;
        }

        h1 {
            color: var(--primary-color);
            text-align: center;
            margin: 20px 0;
            font-size: 2.5em;
            text-shadow: 2px 2px 4px var(--shadow-color);
            position: relative;
            padding-bottom: 15px;
            background: linear-gradient(135deg, var(--primary-color), var(--secondary-color));
            -webkit-background-clip: text;
            -webkit-text-fill-color: transparent;
            font-weight: 800;
            letter-spacing: -0.5px;
        }

        h1::after {
            content: '';
            position: absolute;
            bottom: 0;
            left: 50%;
            transform: translateX(-50%);
            width: 150px;
            height: 4px;
            background: linear-gradient(to right, var(--primary-color), var(--secondary-color));
            border-radius: 2px;
            box-shadow: 0 2px 4px var(--shadow-color);
        }

        .button-container {
            display: flex;
            justify-content: space-between;
            margin-bottom: 20px;
        }

        .back-button {
            background: linear-gradient(135deg, var(--primary-color), var(--primary-dark));
            color: white;
            border: none;
            padding: 12px 24px;
            border-radius: 25px;
            cursor: pointer;
            font-size: 1em;
            font-weight: 600;
            transition: all var(--transition-speed) ease;
            box-shadow: 0 2px 4px var(--shadow-color);
            text-decoration: none;
            display: inline-block;
            position: absolute;
            bottom: 10px;
            left: 10px;
        }

        .back-button:hover {
            transform: translateY(-2px);
            box-shadow: 0 4px 8px var(--shadow-color);
        }

        .back-button:active {
            transform: translateY(0);
        }

        .test-container {
            background: white;
            border-radius: 8px;
            padding: 20px;
            margin: 20px 0;
            box-shadow: 0 2px 4px rgba(0,0,0,0.1);
        }

        .test-form {
            display: flex;
            flex-direction: column;
            gap: 15px;
        }

        .info-section {
            display: flex;
            flex-wrap: wrap;
            gap: 15px;
            margin-bottom: 15px;
            padding: 15px;
            background: #f8f8f8;
            border-radius: 8px;
        }

        .info-item {
            display: flex;
            align-items: center;
            margin-right: 15px;
        }

        .info-label {
            font-weight: bold;
            margin-right: 8px;
            color: var(--text-secondary);
        }

        textarea {
            width: 100%;
            min-height: 200px;
            padding: 10px;
            border: 1px solid #ddd;
            border-radius: 4px;
            font-family: monospace;
            font-size: 14px;
            resize: vertical;
        }

        .submit-button {
            background: linear-gradient(135deg, var(--primary-color), var(--primary-dark));
            color: white;
            border: none;
            padding: 12px 24px;
            border-radius: 25px;
            cursor: pointer;
            font-size: 1em;
            font-weight: 600;
            transition: all var(--transition-speed) ease;
            box-shadow: 0 2px 4px var(--shadow-color);
            align-self: flex-start;
        }

        .submit-button:hover {
            transform: translateY(-2px);
            box-shadow: 0 4px 8px var(--shadow-color);
        }

        .submit-button:active {
            transform: translateY(0);
        }

        .result-container {
            margin-top: 20px;
            padding: 20px;
            background: #f9f9f9;
            border-radius: 8px;
            border-left: 4px solid var(--primary-color);
            display: none;
        }

        .result-title {
            font-weight: bold;
            margin-bottom: 10px;
            color: var(--text-primary);
            font-size: 1.2em;
        }

        .result-item {
            margin: 10px 0;
            display: flex;
            align-items: center;
        }

        .result-label {
            font-weight: bold;
            min-width: 150px;
            color: var(--text-secondary);
        }

        .footer {
            position: fixed;
            bottom: 0;
            left: 0;
            right: 0;
            text-align: center;
            padding: 10px;
            background: rgba(255, 255, 255, 0.9);
            backdrop-filter: blur(10px);
            z-index: 100;
        }

        .footer-image {
            max-width: 100%;
            height: auto;
            display: block;
            margin: 0 auto;
            width: 300px;
        }

        /* Badge styles */
        .icon-badge {
            display: inline-flex;
            align-items: center;
            padding: 4px 10px;
            border-radius: 12px;
            font-size: 0.85em;
            font-weight: 600;
            gap: 4px;
            box-shadow: 0 1px 2px rgba(0,0,0,0.1);
        }

        .icon-badge-text {
            margin-left: 4px;
        }

        .icon-badge-true {
            background-color: var(--badge-success-bg);
            color: var(--badge-success-text);
            border: 1px solid rgba(30, 142, 62, 0.2);
        }

        .icon-badge-false {
            background-color: var(--badge-error-bg);
            color: var(--badge-error-text);
            border: 1px solid rgba(217, 48, 37, 0.2);
        }

        .icon-badge-backward {
            background-color: var(--badge-warning-bg);
            color: var(--badge-warning-text);
            border: 1px solid rgba(245, 124, 0, 0.2);
        }

        .icon-badge-forward {
            background-color: var(--badge-info-bg);
            color: var(--badge-info-text);
            border: 1px solid rgba(26, 115, 232, 0.2);
        }

        .icon-badge-full {
            background-color: var(--badge-type-bg);
            color: var(--badge-type-text);
            border: 1px solid rgba(156, 39, 176, 0.2);
        }

        .icon-badge-type {
            background-color: var(--badge-type-bg);
            color: var(--badge-type-text);
            border: 1px solid rgba(156, 39, 176, 0.2);
        }

        .icon-badge-none {
            background-color: var(--badge-neutral-bg);
            color: var(--badge-neutral-text);
            border: 1px solid rgba(95, 99, 104, 0.2);
        }

        .icon-badge-warning {
            background-color: var(--badge-warning-bg);
            color: var(--badge-warning-text);
            border: 1px solid rgba(245, 124, 0, 0.2);
        }

        .icon-badge-subject {
            background-color: var(--badge-info-bg);
            color: var(--badge-info-text);
            border: 1px solid rgba(26, 115, 232, 0.2);
        }

        .icon-badge-version {
            background-color: var(--badge-type-bg);
            color: var(--badge-type-text);
            border: 1px solid rgba(156, 39, 176, 0.2);
        }

        .icon-badge-id {
            background-color: var(--badge-neutral-bg);
            color: var(--badge-neutral-text);
            border: 1px solid rgba(95, 99, 104, 0.2);
        }
    </style>
</head>
<body>
    <div class="header-container">
        <a href="/schema/?topic={{.SubjectName}}" class="back-button">Back to Schema View</a>
        <img src="/static/header.png" class="header-image" alt="Header">
        <div class="header-stats">
            <a href="https://slack.com" target="_blank" class="slack-button">
                <img src="/static/slack-channel.png" class="slack-image" alt="Slack">
            </a>
            <a href="https://www.lemonde.fr" target="_blank" class="github-button">
                <img src="/static/github.png" class="github-image" alt="GitHub">
            </a>
        </div>
        <h1>Test Schema Compatibility</h1>
    </div>

    <div class="test-container">
        <div class="info-section">
            <div class="info-item">
                <span class="info-label">Subject:</span>
                <span class="icon-badge icon-badge-subject">{{.SubjectName}}</span>
            </div>
            <div class="info-item">
                <span class="info-label">Version:</span>
                <span class="icon-badge icon-badge-version">{{.Version}}</span>
            </div>
            <div class="info-item">
                <span class="info-label">ID:</span>
                <span class="icon-badge icon-badge-id">{{.SchemaID}}</span>
            </div>
        </div>

        <div class="test-form">
            <label for="testJson">Enter JSON to test compatibility:</label>
            <textarea id="testJson" placeholder="Paste your JSON here..."></textarea>
            <button id="testButton" class="submit-button">Test Compatibility</button>
        </div>

        <div id="resultContainer" class="result-container">
            <div class="result-title">Compatibility Test Results</div>
            <div class="result-item">
                <span class="result-label">Compatibility:</span>
                <span id="compatibilityResult"></span>
            </div>
            <div class="result-item">
                <span class="result-label">HTTP Status:</span>
                <span id="statusResult"></span>
            </div>
            <div class="result-item">
                <span class="result-label">Error Code:</span>
                <span id="errorCodeResult"></span>
            </div>
            <div class="result-item">
                <span class="result-label">Message:</span>
                <span id="messageResult"></span>
            </div>
        </div>
    </div>

    <div class="footer">
        <img src="/static/footer.png" class="footer-image" alt="Footer">
    </div>

    <script>
        document.getElementById('testButton').addEventListener('click', testSchema);

        function testSchema() {
            const testJsonText = document.getElementById('testJson').value;
            const subject = "{{.SubjectName}}";
            const version = "{{.Version}}";
            const id = "{{.SchemaID}}";

            // Show loading state
            const testButton = document.getElementById('testButton');
            const originalButtonText = testButton.textContent;
            testButton.textContent = 'Testing...';
            testButton.disabled = true;

            fetch('/test-schema', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify({
                    subject: subject,
                    version: version,
                    id: id,
                    json: testJsonText
                })
            })
            .then(response => response.json())
            .then(data => {
                // Reset button
                testButton.textContent = originalButtonText;
                testButton.disabled = false;

                // Determine badge classes based on response values
                let compatibilityBadgeClass = "icon-badge-none";
                if (data.isCompatible === true) {
                    compatibilityBadgeClass = "icon-badge-true";
                } else if (data.isCompatible === false) {
                    compatibilityBadgeClass = "icon-badge-false";
                } else {
                    compatibilityBadgeClass = "icon-badge-warning";
                }

                let statusBadgeClass = "icon-badge-none";
                if (data.httpStatus >= 400) {
                    statusBadgeClass = "icon-badge-false";
                } else if (data.httpStatus >= 300) {
                    statusBadgeClass = "icon-badge-warning";
                } else if (data.httpStatus >= 200) {
                    statusBadgeClass = "icon-badge-true";
                }

                let errorBadgeClass = "icon-badge-none";
                if (data.errorCode > 0) {
                    errorBadgeClass = "icon-badge-false";
                }

                let messageBadgeClass = "icon-badge-none";
                if (data.message && data.message !== "None") {
                    if (data.isCompatible === false) {
                        messageBadgeClass = "icon-badge-false";
                    } else if (data.isCompatible === true) {
                        messageBadgeClass = "icon-badge-true";
                    }
                }

                // Update the result display
                document.getElementById('compatibilityResult').innerHTML = 
                    '<span class="icon-badge ' + compatibilityBadgeClass + '">' + (data.isCompatible === true ? 'Compatible' : data.isCompatible === false ? 'Not Compatible' : 'Unknown') + '</span>';
                
                document.getElementById('statusResult').innerHTML = 
                    '<span class="icon-badge ' + statusBadgeClass + '">' + data.httpStatus + '</span>';
                
                document.getElementById('errorCodeResult').innerHTML = 
                    '<span class="icon-badge ' + errorBadgeClass + '">' + (data.errorCode === 0 ? 'None' : data.errorCode) + '</span>';
                
                document.getElementById('messageResult').innerHTML = 
                    '<span class="icon-badge ' + messageBadgeClass + '">' + (data.message && data.message !== "None" ? data.message : 'None') + '</span>';

                // Show the result container
                document.getElementById('resultContainer').style.display = 'block';
            })
            .catch(error => {
                // Reset button
                testButton.textContent = originalButtonText;
                testButton.disabled = false;
                
                // Show error in the result container
                document.getElementById('compatibilityResult').innerHTML = 
                    '<span class="icon-badge icon-badge-false">Error</span>';
                
                document.getElementById('statusResult').innerHTML = 
                    '<span class="icon-badge icon-badge-false">Error</span>';
                
                document.getElementById('errorCodeResult').innerHTML = 
                    '<span class="icon-badge icon-badge-false">API Error</span>';
                
                document.getElementById('messageResult').innerHTML = 
                    '<span class="icon-badge icon-badge-false">' + (error.message || 'Failed to test schema') + '</span>';
                
                // Show the result container
                document.getElementById('resultContainer').style.display = 'block';
            });
        }
    </script>
</body>
</html>`
