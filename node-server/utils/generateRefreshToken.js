import jwt from "jsonwebtoken";

export default function generateRefreshToken(res, data) {
  return jwt.sign(data, process.env.REFRESH_TOKEN_SECRET, {
    expiresIn: "1d",
  });
}
