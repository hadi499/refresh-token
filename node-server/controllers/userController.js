import asyncHandler from "express-async-handler";
import User from "../models/userModel.js";
import generateAccessToken from "../utils/generateAccessToken.js";
import generateRefreshToken from "../utils/generateRefreshToken.js";
import jwt from "jsonwebtoken";

// Register new user
const registerUser = asyncHandler(async (req, res) => {
  const { name, email, password } = req.body;

  const userExists = await User.findOne({ email });
  if (userExists) {
    res.status(400);
    throw new Error("User already exists");
  }

  const user = await User.create({ name, email, password });

  if (user) {
    res.status(201).json({ message: "User registered successfully" });
  } else {
    res.status(400);
    throw new Error("Invalid user data");
  }
});

// Login user
const authUser = asyncHandler(async (req, res) => {
  const { email, password } = req.body;

  const user = await User.findOne({ email });

  if (user && (await user.matchPassword(password))) {
    const data = {
      id: user.id,
      name: user.name,     
      email: user.email,
    };
    const accessToken = generateAccessToken(res, data);
    const refreshToken = generateRefreshToken(res, data);

    user.refreshToken = refreshToken;
    await user.save();

    res.json({ accessToken, refreshToken });
  } else {
    res.status(401);
    throw new Error("Invalid email or password");
  }
});

// Get all users (for testing)
const getUser = asyncHandler(async (req, res) => {
  const users = await User.find();
  if (!users) {
    res.status(404);
    throw new Error("No users found");
  }
  res.json(users);
});

// Refresh access token using refresh token
const myRefreshToken = asyncHandler(async (req, res) => {
  const { token: refreshToken } = req.body;

  if (!refreshToken) return res.sendStatus(401);

  const user = await User.findOne({ refreshToken });
  if (!user) return res.sendStatus(403);

  jwt.verify(refreshToken, process.env.REFRESH_TOKEN_SECRET, (err, decoded) => {
    if (err) return res.sendStatus(403);

    const data = { id: decoded.id, email: decoded.email };
    const newAccessToken = generateAccessToken(res, data);

    res.json({ accessToken: newAccessToken, refreshToken });
  });
});

// Logout user (invalidate refresh token)
const logoutUser = asyncHandler(async (req, res) => {
  const { token: refreshToken } = req.body;

  if (!refreshToken) return res.sendStatus(400);

  const user = await User.findOne({ refreshToken });

  if (!user) return res.sendStatus(204); // No content, token already removed

  user.refreshToken = "";
  await user.save();

  res.json({ message: "Logged out successfully" });
});

export { registerUser, authUser, getUser, myRefreshToken, logoutUser };
