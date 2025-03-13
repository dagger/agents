package io.dagger.modules.java;

import io.dagger.client.*;
import io.dagger.module.AbstractModule;
import io.dagger.module.annotation.DefaultPath;
import io.dagger.module.annotation.Function;
import io.dagger.module.annotation.Object;
import java.util.List;
import java.util.concurrent.ExecutionException;

@Object
public class Java extends AbstractModule {
  public Directory source;

  public Java() {
    super();
  }

  /**
   * @param source Directory containing the sources to work on
   */
  public Java(Client dag, @DefaultPath(".") Directory source) {
    super(dag);
    this.source = source;
  }

  /** Let AI find and explain potential bugs in the source code. */
  @Function
  public String findBugs() throws ExecutionException, DaggerQueryException, InterruptedException {
    var llm =
        asReader(
            "find potential bugs in the existing java code, explain them and propose alternative code to fix them");
    return llm.lastReply();
  }

  /** Refactor existing code to improve readability and maintainability. */
  @Function
  public Directory refactor() {
    var llm =
        asEditor("refactor the existing java code to improve readability and maintainability");
    return llm.Workspace().dir();
  }

  /**
   * Ask anything to the AI in the context of an editor.
   *
   * <p>Files will be edited in place and a {@code .nono.md} file will contain some explanation
   *
   * @param assignment Assignment for the AI
   */
  @Function
  public Llm asEditor(String assignment) {
    return askTo("editor", assignment);
  }

  /**
   * Ask anything to the AI in the context of a reader without touching any file
   *
   * @param assignment Assignment for the AI
   */
  @Function
  public Llm asReader(String assignment) {
    return askTo("reader", assignment);
  }

  public Llm askTo(String profile, String assignment) {
    return dag.llm()
        .withWorkspace(
            dag.workspace(
                new Client.WorkspaceArguments()
                    .withStart(this.source)
                    .withChecker(
                        mvnContainer().withDefaultArgs(List.of("mvn", "test", "compile")))))
        .withPromptVar("assignment", assignment)
        .withPromptFile(
            dag.currentModule().source().file("src/main/resources/prompts/" + profile + ".txt"));
  }

  public Container mvnContainer() {
    return dag.container()
        .from("maven:3.9.9-eclipse-temurin-17")
        .withMountedCache("/root/.m2", dag.cacheVolume("m2_cache"))
        .withWorkdir("/app");
  }
}
