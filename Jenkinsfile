pipeline {
    agent {
        kubernetes {
            label 'jenkins-agent'
            defaultContainer 'jnlp'
            yaml """
apiVersion: v1
kind: Pod
metadata:
  labels:
    some-label: jenkins-agent
spec:
  containers:
  - name: go-builder
    image: golang:1.18
    command:
    - cat
    tty: true
    resources:
      limits:
        cpu: "500m"
        memory: "512Mi"
      requests:
        cpu: "250m"
        memory: "256Mi"
"""
        }
    }
    environment {
        DOCKER_IMAGE = "your-docker-repo/your-app"
        DOCKER_TAG = "latest"
    }
    stages {
        stage('Checkout') {
            steps {
                git 'https://github.com/your-repo/your-go-k8s-project.git'
            }
        }
        stage('Build') {
            steps {
                container('go-builder') {
                    sh 'go build -o your-app'
                }
            }
        }
        stage('Test') {
            steps {
                container('go-builder') {
                    sh 'go test ./...'
                }
            }
        }
        stage('Docker Build & Push') {
            steps {
                script {
                    // 构建镜像并推送到 Docker 仓库
                    def appImage = docker.build("${DOCKER_IMAGE}:${DOCKER_TAG}")
                    docker.withRegistry('https://your-docker-registry', 'docker-credentials-id') {
                        appImage.push()
                    }
                }
            }
        }
        stage('Deploy to Kubernetes') {
            steps {
                script {
                    // 调用 kubectl 或使用 Kubernetes Deploy 插件进行部署
                    sh 'kubectl apply -f k8s/deployment.yaml'
                }
            }
        }
    }
}
