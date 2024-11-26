# Gator CLI - Blog Aggregator

Gator CLI is a command-line tool for aggregating and managing RSS feeds. Users can log in, follow feeds, and browse posts from their favorite blogs.

## Requirements

Before running the Gator CLI, ensure you have the following installed on your system:

1. **PostgreSQL**: The database used to store users, feeds, and posts.
2. **Go**: The programming language used to develop the Gator CLI.

### Install PostgreSQL
Follow the [official PostgreSQL installation guide](https://www.postgresql.org/download/) for your operating system. Note your connection details (username, password, host, port, and database name) as they will be needed for configuration.

### Install Go
Follow the [official Go installation guide](https://go.dev/doc/install) to install the latest version of Go.

---

## Installation

To install the Gator CLI, use the following command in your terminal:

```bash
go install github.com/seanhuebl/blog_aggregator/cmd/gator@latest
```

This command will fetch, build, and install the CLI tool into your Go binary directory. Ensure the Go binary path is added to your system's PATH to use `gator` globally.

---

## Configuration

To run the Gator CLI, you need to set up a configuration file. This file should be created in your home directory with the name `.gatorconfig.json`. 

### Example Configuration File

```json
{
  "db_url": "connection_string_goes_here",
  "current_user_name": "username_goes_here"
}
```

### Database Connection String Format

The `db_url` field must follow this format:
```
protocol://username:password@host:port/database?sslmode=disable
```

- Replace `protocol`, `username`, `password`, `host`, `port`, and `database` with your PostgreSQL setup details.
- Use `sslmode=disable` unless SSL is explicitly required for your database.

### Example Connection String

```json
"db_url": "postgres://user:password@localhost:5432/mydatabase?sslmode=disable"
```

---

## Running the Program

Once the configuration file is set up, you can run the Gator CLI using:

```bash
gator <command> [arguments]
```

### Available Commands

1. **Login**: Set the current user for the CLI.
   ```bash
   gator login <username>
   ```

2. **Register**: Create a new user in the system.
   ```bash
   gator register <username>
   ```

3. **AddFeed**: Add a new feed and follow it.
   ```bash
   gator addfeed <feed_name> <feed_url>
   ```

4. **Follow**: Follow an existing feed.
   ```bash
   gator follow <feed_url>
   ```

5. **Unfollow**: Unfollow a feed.
   ```bash
   gator unfollow <feed_url>
   ```

6. **Feeds**: List all available feeds.
   ```bash
   gator feeds
   ```

7. **Following**: List feeds you are currently following.
   ```bash
   gator following
   ```

8. **Browse**: Browse posts from feeds you follow. Optionally specify the number of posts to retrieve.
   ```bash
   gator browse [limit]
   ```

9. **Reset**: Delete all users and reset the system.
   ```bash
   gator reset
   ```

10. **Aggregate (Agg)**: Periodically fetch new posts from feeds.
    ```bash
    gator agg <interval>
    ```
    - `<interval>`: Time duration between fetches (e.g., `30s`, `5m`, `1h`).

---

## Example Workflow

1. **Set up PostgreSQL and configure the connection string** in `.gatorconfig.json`.
2. **Run the CLI**:
   - Register a user: `gator register user1`
   - Add and follow a feed: `gator addfeed "Tech News" "https://example.com/rss"`
   - In a seperate terminal run `gator agg <interval>`
   - Browse posts: `gator browse 5`

---
## My Learning Journey
This project was part of a guided learning experience with minimal pseudocode provided. It introduced me to several important concepts in building backend systems and working with databases, including:

**Connecting to Databases**: Using Go and PostgreSQL to manage data storage and retrieval.
**Managing Migrations**: Using the goose library to handle database schema migrations seamlessly.
**Generating Go Code for Queries**: Leveraging sqlc to generate type-safe Go code directly from SQL queries.
Through this project, I’ve gained hands-on experience in integrating these tools and begun to understand the fundamentals of how full web applications are developed. It's an exciting step toward building more complex and feature-rich applications in the future!

---
## Troubleshooting

- Ensure PostgreSQL is running and accessible with the provided connection string.
- Verify the configuration file exists and is correctly formatted.
- Check that the `gator` binary is installed and accessible from your system PATH.

---

## License

This project is licensed under the [MIT License](LICENSE).

---

Enjoy using Gator CLI!