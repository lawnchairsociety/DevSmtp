namespace DevSmtp.Core.Commands
{
    public class CommandResult
    {
        public CommandResult()
        {
            this.Succeeded = true;
        }

        public CommandResult(Exception error)
        {
            this.Error = error;
            this.Succeeded = false;
        }

        public bool Succeeded { get; }
        public Exception? Error { get; }
    }
}
