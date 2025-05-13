import express from "express";

import {
  registerUser,
  authUser,
  getUser,
  myRefreshToken,
  logoutUser
} from "../controllers/userController.js";
import authenticateToken from "../middleware/authenticateToken.js";

const router = express.Router();

router.post("/register", registerUser);
router.post("/auth", authUser);
router.get("/", authenticateToken, getUser);
router.post("/logout",  logoutUser);
router.post("/refresh", myRefreshToken);


export default router;
