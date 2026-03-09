# World ID Setup

This project uses **World ID** to verify that the user is a real human.

This prevents **Sybil attacks** in the tokenization process.

Official website:

https://world.org/

Developer documentation:

https://docs.world.org/

Developer portal:

https://developer.worldcoin.org/

---

# 1. Install World App

Download the mobile application:

https://world.org/

Click **Get World App**.

Install it on your phone.

---

# 2. Create an Account

Open the app and create your account.

Make sure to **protect the account with an access key**.

This will be used later for developer login.

---

# 3. Access the Developer Portal

Go to:

https://developer.worldcoin.org/login

Click **Continue with World ID**.

Use the World App on your phone to authenticate.

---

# 4. Create Your First Team

Inside the developer portal:

Create a **Team**.

Example name:

```
hackthonteam
```

You can keep the default configuration.

---

# 5. Create a World ID App

Inside your team:

Create a new **Application**.

Suggested name:

```
autolocktokenization
```

---

# 6. Enable World ID 4.0 (Managed)

Inside the app configuration:

Enable:

```
World ID 4.0 Managed
```

This allows Worldcoin to handle the verification infrastructure.

---

# 7. Retrieve Configuration Values

Go to the **World ID 4.0 tab** in the developer portal.

Copy the following values:

```
app_id
rp_id
action_name
```

These values will be used in:

```
auto-lock-defi/config.staging.json
```
```
frontend/.env.local
```
---

# 8. Generate an API Key

Inside the developer portal:

Generate a **new API key**.

Save this value.

You will use it in the environment variables.

Copy and paste it for both env vars below inside the `.env.local` in the `frontend` project.

Example:

```
NEXT_PUBLIC_WORLD_ID_API_KEY_ALL="0xfcf93..."
RP_SIGNING_KEY="0xfcf93..."
```

---

# 9. Create an Action

Navigate to the **Actions** tab.

Create your first action.

Example:

```
tokenizevehicle
```

This action will be referenced when verifying proofs.

---

# Setup Completed

Your World ID integration is now ready.

Continue with:

* ⛓️ [Tenderly Setup](README_TENDERLY.md)

