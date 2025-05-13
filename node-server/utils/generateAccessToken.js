import jwt from "jsonwebtoken";

export default function generateAccessToken(res, data) {
  return jwt.sign(data, process.env.ACCESS_TOKEN_SECRET, {
    expiresIn: "3m",
  });
}
