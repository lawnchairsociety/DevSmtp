namespace DevSmtp.Core.Commands
{
    public sealed class HeloResult : CommandResult
    {
        public HeloResult()
        {
        }

        public HeloResult(Exception error)
            : base(error)
        {
        }
    }
}
