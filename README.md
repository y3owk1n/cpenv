# cpenv

Copy and paste your local .env to right projects faster

## Usage

Install cpenv globally:

```bash
npm install -g cpenv
```

Make a local env vault at `~/env-files`, and restructure the folders with env as per your actual project.

```bash
- ~/env-files
├── single-env-project
├──── .env
├── multi-env-project
├──── .env
├──── apps
├────── web
├──────── .env
├────── api
├──────── .env
```

Go to your project directory and run the following command in your terminal:

```bash
cpenv
```
