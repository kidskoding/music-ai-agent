# TODO.md — Music Agent Project

## Core Agent

- [x] Store tracks and session memory (history, skips, energy, mood).  
- [x] Decide the next track based on current mood, energy, and rules.  
- [x] Avoid recently played tracks and handle fallback when no matches are found.  
- [ ] Randomize selections among candidate tracks to create a natural listening flow.  

---

## Testing

- [x] Unit tests for track selection logic.  
- [x] Sanity tests for framework verification.  
- [x] Tests for skip/recently played behavior and fallback logic.  
- [x] Tests for energy/mood adaptive behavior.  
- [ ] Table-driven tests for multiple session scenarios.  
- [ ] Integration tests with Databricks and Google Cloud AI (mocked or real).  

---

## Agent Behavior Enhancements

- [ ] Adaptive mood and energy transitions based on session memory and AI input.  
- [ ] Skip history management to prevent repeated skips.  
- [ ] Genre weighting and preference support.  
- [ ] Persistent session memory across runs and devices.  

---

## Databricks Integration

- [ ] Store and retrieve session memory for analytics.  
- [ ] Track user listening patterns and energy transitions.  
- [ ] Generate recommendation signals from processed data.  
- [ ] Persist memory in cloud tables or Delta tables for cross-session usage.  
- [ ] Dynamic track queries from Databricks instead of static samples.  

---

## Google Cloud AI Integration

- [ ] Personalized recommendations using Vertex AI or other GCP services.  
- [ ] Next-track prediction using AI models.  
- [ ] Dynamic agent behavior updates based on AI feedback.  
- [ ] Logging AI model decisions and confidence scores for observability.  

---

## Tooling and Workflow

- [x] Makefile for running tests, building, and running the agent.  
- [ ] Modular project structure: agent logic, memory, Spotify integration, logger, utilities.  
- [ ] Configuration management for Databricks and Google Cloud credentials.  

---

## Logging and Observability

- [ ] Log agent decisions (track, mood, energy, AI recommendation).  
- [ ] Track session memory changes.  
- [ ] Enable debug or verbose mode for troubleshooting.  

---

## Future / Optional Enhancements

- [ ] Multi-agent orchestration for moods, genres, or energy levels.  
- [ ] Real-time playlist updates based on agent decisions.  
- [ ] Feedback loop for likes, skips, and user preferences.  
- [ ] Integration of machine learning for next-track prediction.  
- [ ] Visualization of agent decisions and session history in a web or desktop UI.
