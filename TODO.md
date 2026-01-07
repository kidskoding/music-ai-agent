# TODO.md — Music AI Agent

## Core Agent

- [x] Store tracks and session memory (history, skips, energy, mood).  
- [x] Decide the next track based on current mood, energy, and rules.  
- [x] Avoid recently played tracks and handle fallback when no matches are found.  
- [x] Randomize selections among candidate tracks to create a natural listening flow.  

---

## Testing

- [x] Unit tests for track selection logic.  
- [x] Sanity tests for framework verification.  
- [x] Tests for skip/recently played behavior and fallback logic.  
- [x] Tests for energy/mood adaptive behavior.  
- [x] Integration tests for event logging (LocalStore verification).
- [x] Edge case testing (Empty lists, unknown moods, all tracks skipped).  
- [x] Integration tests with Databricks and Google Cloud AI (Real connection verified).  

---

## Agent Behavior Enhancements

- [ ] Adaptive mood and energy transitions based on session memory and AI input.  
- [ ] Skip history management to prevent repeated skips.  
- [ ] Genre weighting and preference support.  
- [ ] Persistent session memory across runs and devices.  

---

## Spotify Integration

- [x] Implement OAuth2 Client for Spotify authentication.
- [ ] Replace static track list with live Spotify `UserTopTracks`.
- [ ] Implement `QueueTrack` to control real playback.
- [ ] Handle API rate limiting and token refreshing.

---

## Databricks Integration

- [x] Store and retrieve session memory for analytics.  
- [x] Track user listening patterns and energy transitions.  
- [x] Generate recommendation signals from processed data.  
- [x] Persist memory in cloud tables or Delta tables for cross-session usage.  
- [x] Dynamic track queries from Store/DB instead of static samples.  

---

## Google Cloud AI Integration

- [x] Personalized recommendations using Vertex AI
- [x] Parse AI response to influence track selection logic.
- [x] Next-track prediction using AI models.  
- [ ] Dynamic agent behavior updates based on AI feedback.  
- [ ] Logging AI model decisions and confidence scores for observability.  

---

## Tooling and Workflow

- [x] Makefile for running tests, building, and running the agent.  
- [x] Modular project structure: agent logic, memory, logger, utilities.  
- [x] Refactored `Track` model to shared package to prevent circular imports.
- [x] Configuration management for Databricks, Spotify, and Google Cloud credentials.  

---

## Logging and Observability

- [x] Log agent decisions (track, mood, energy, AI recommendation).  
- [x] Track session memory changes.  
- [ ] Enable debug or verbose mode for troubleshooting.  

---

## Future / Optional Enhancements

- [ ] Multi-agent orchestration for moods, genres, or energy levels.  
- [ ] Real-time playlist updates based on agent decisions.  
- [ ] Feedback loop for likes, skips, and user preferences.  
- [ ] Integration of machine learning for next-track prediction.  
- [ ] Visualization of agent decisions and session history in a web or desktop UI.
