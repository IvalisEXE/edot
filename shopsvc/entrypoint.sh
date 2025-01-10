#!/bin/sh

# Tunggu PostgreSQL siap (menggunakan pg_isready)
echo "Menunggu PostgreSQL..."
until pg_isready -h mydb -p 5432; do
  sleep 1
done

echo "PostgreSQL siap!"

# Tunggu Redis siap
echo "Menunggu Redis..."
until redis-cli -h myredis ping; do
  sleep 1
done


# Run migration
echo "Running migrations..."
./api migrate

# Check if migration was successful
if [ $? -ne 0 ]; then
  echo "Migration failed. Exiting..."
  exit 1
fi

# Run seeder
echo "Seeding database..."
./api seed

# Check if seeder was successful
if [ $? -ne 0 ]; then
  echo "Seeder failed. Exiting..."
  exit 1
fi