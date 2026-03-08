# Tenderly Setup

This project uses **Tenderly Virtual Testnets** as the blockchain environment.

Official website:

https://tenderly.co/

---

# 1. Create an Account

Go to:

https://tenderly.co/

Click **Go to Dashboard**.

Login or register.

You can use:

* Google
* GitHub

---

# 2. Create an Organization

Navigate to:

```
Members → Create Organization
```

Use the **Free Trial plan**.

Example name:

```
AutoLock
```

---

# 3. Create a Project

Inside your organization:

Create a new project.

Example name:

```
AutoLock DeFi
```

---

# 4. Create a Virtual Testnet

Inside the project dashboard:

Go to:

```
Virtual Testnets
```

Click:

```
Create Virtual Testnet
```

---

# 5. Copy RPC URLs

After creating the testnet you will see:

```
Admin RPC (HTTPS)
Admin RPC (WebSocket)
```

Copy both values.

---

# 6. Configure Environment Variables

Add the RPC URL to your `.env` file.

Example:

```
TENDERLY_RPC_URL=https://rpc.tenderly.co/...
```

---

# 7. Configure CRE Project

Open:

```
project.yaml
```

Set the staging RPC URL using the same value.

Example:

```
rpc_url: https://rpc.tenderly.co/...
```

---

# Setup Completed

Tenderly is now ready.

You can now run the project deployment.
