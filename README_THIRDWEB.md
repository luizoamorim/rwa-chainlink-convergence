# Thirdweb Setup

This project uses **Thirdweb** to connect user wallets and interact with the blockchain from the frontend.

Official website:

https://thirdweb.com/

---

# 1. Create an Account

Go to:

https://thirdweb.com/

Click **Start for free**.

Sign up using **Google**.

---

# 2. Choose the Trial Plan

Use the **Free Trial Account**.

For this project it is enough.

---

# 3. Skip Team Members

When prompted to add team members:

Skip this step for now.

---

# 4. Create a Project

Create a new project.

You can use any name.

Example:

```
AutoLock DeFi
```

---

# 5. Configure Allowed Domains

Allow **all domains** for now to simplify testing.

Later this can be restricted.

---

# 6. Save Your Client ID

Inside your project dashboard you will find the **Client ID**.

Copy this value.

---

# 7. Configure Frontend Environment

Open:

```
frontend/.env.local
```

Add your client id:

```
NEXT_PUBLIC_THIRDWEB_CLIENT_ID=your_client_id_here
```

---

# Setup Completed

Your Thirdweb integration is now ready.

You can continue with the next setup:

* 🧑‍🚀 [World ID Setup](README_WORLD_ID.md)
* ⛓️ [Tenderly Setup](README_TENDERLY.md)
