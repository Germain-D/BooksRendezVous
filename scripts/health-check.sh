#!/bin/bash

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo "🏥 Books Rendezvous Health Check"
echo "================================"

# Check if Docker Compose is running
if ! docker-compose ps | grep -q "Up"; then
    echo -e "${RED}❌ Docker Compose services are not running${NC}"
    echo "Run 'docker-compose up -d' to start the services"
    exit 1
fi

# Check Frontend
echo -n "🌐 Frontend (http://localhost:3000): "
if curl -f -s http://localhost:3000 > /dev/null; then
    echo -e "${GREEN}✅ Healthy${NC}"
else
    echo -e "${RED}❌ Not responding${NC}"
fi

# Check Backend
echo -n "🔧 Backend (http://localhost:6050): "
if curl -f -s http://localhost:6050/health > /dev/null; then
    echo -e "${GREEN}✅ Healthy${NC}"
else
    echo -e "${RED}❌ Not responding${NC}"
fi

# Check Database
echo -n "🗄️  Database (localhost:5432): "
if docker-compose exec -T database pg_isready -U bibliotheque_user -d bibliotheque > /dev/null 2>&1; then
    echo -e "${GREEN}✅ Healthy${NC}"
else
    echo -e "${RED}❌ Not responding${NC}"
fi

# Check Redis (optional)
echo -n "🔴 Redis (localhost:6379): "
if docker-compose exec -T redis redis-cli ping > /dev/null 2>&1; then
    echo -e "${GREEN}✅ Healthy${NC}"
else
    echo -e "${YELLOW}⚠️  Not responding (optional service)${NC}"
fi

echo ""
echo "📊 Service Status:"
docker-compose ps 