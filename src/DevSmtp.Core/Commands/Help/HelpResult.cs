namespace DevSmtp.Core.Commands
{
    public sealed class HelpResult : CommandResult
    {
        public HelpResult()
        {
        }

        public HelpResult(Exception error)
            : base(error)
        {
        }
    }
}
