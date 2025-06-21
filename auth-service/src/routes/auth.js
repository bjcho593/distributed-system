import express from 'express';
import jwt from 'jsonwebtoken';
import { redisClient } from '../services/redisClient.js';

const router = express.Router();

router.post('/login', async (req, res) => {
  const { username, password } = req.body;

  // Validación dummy
  if (username === 'admin' && password === 'admin') {
    const token = jwt.sign({ user: username }, process.env.JWT_SECRET, { expiresIn: '1h' });

    // Guardar en Redis
    await redisClient.setEx(token, 3600, username);

    return res.json({ token });
  }

  res.status(401).json({ message: 'Invalid credentials' });
});

export default router;
